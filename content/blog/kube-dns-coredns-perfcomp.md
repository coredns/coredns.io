+++
date = "2018-05-31T06:21:00Z"
description = "Relative Performance of CoreDNS vs Kube-dns"
tags = ["Kubernetes", "Performance", "Kube-DNS", "CoreDNS"]
title = "Performance Testing CoreDNS and Kube-dns"
author = "chris"
+++

# Relative Performance of CoreDNS vs Kube-dns

CoreDNS is on track to replace kube-dns as the the default cluster DNS service in Kubernetes.  Since Kubernetes 1.9, CoreDNS has been available as an optional replacement, and now CoreDNS is proposed to replace kube-dns as the default cluster DNS service in 1.12. As part of this transition, we have extended the existing kube-dns oriented [performance testing tool](https://github.com/kubernetes/perf-tests/tree/master/dns) to also [work with CoreDNS](https://github.com/kubernetes/perf-tests/pull/114).  This blog covers a performance comparison of kube-dns and CoreDNS using the updated tool.

## The Kubernetes DNS Perf-Test Tool

The [Kubernetes performance testing tool](https://github.com/kubernetes/perf-tests/tree/master/dns) is an set of python scripts that automate testing various configurations of kube-dns and CoreDNS in a small cluster.  It is geared toward low scale performance testing with the DNS servers under test running in CPU restricted containers. The tool automates configuring the DNS servers, starting them up, spinning the test client up, and running tests.  The output produced can be digested into a sqlite database using an script in the same repository.  The configurable test options are ...

* Server container CPU allowance
* Server cache size
* Client max QPS
* Client queries - e.g. what names to query
* Test run time (per configuration)

## CoreDNS and Kube-dns Differences

I won't go into a comprehensive list of differences between CoreDNS and kube-dns here because most of them are not related to performance. But I’ll point out a few key differences that do weigh in on performance.

###C vs Go

CoreDNS is a single process, written in Go. Kube-dns actually operates as a few separate processes, the two we care about are *dnsmasq* which handles caching and routing upstream queries, and the namesake *kube-dns* which handles in-cluster names.  The *kube-dns* process is written in Go, but *dnsmasq* is written in C.  As it relates to performance, dnsmasq’s C code should be faster than CoreDNS’s Go implementation at caching and forwarding requests upstream. 

###Multithreading

Dnsmasq is single threaded, so its CPU usage is restricted to one core. CoreDNS is multi-threaded, to it can scale higher in a single instance by consuming multiple cores.  In kube-dns, all incoming DNS requests are routed through dnsmasq.  This presents a bottleneck which can be solved with horizontal scaling (i.e. load balanced replicas of kube-dns).  This horizontal scaling results in a larger memory footprint. It also complicates apples to apples testing a bit.

###Maximum Concurrent Queries

In kube-dns deployments, dnsmasq defaults to allowing a maximum of 150 concurrent queries.  If this number is exceeded, then dnsmasq will drop queries.  This can be a problem in a large scale deployment where you might see DNS request spikes (e.g. hundreds of pods hitting the DNS at nearly exactly the same time).  This can be adjusted up for environments that see big spikes in concurrent DNS requests.  CoreDNS does not have a similar hard limit for the number of concurrent queries.  In theory, instead of dropping packets, CoreDNS should handle big spikes, but with a higher latency.

###Negative Caching

Another performance related difference is how the two handle negative caching. Negative caching is a [standard](https://tools.ietf.org/html/rfc2308) that enables dns servers to cache responses that have no records.  In a DNS response, each record in the response will have a TTL.  A DNS cache uses this TTL to determine how long to keep the record in cache.  Empty responses, and denial of existence responses don’t normally have any records in them, so there is no TTL.  The standard way to handle this is to include an SOA record in the response with a TTL, so downstream caches know how long to keep the response.  In Kubernetes default deployments, CoreDNS keeps a negative cache, whereas [negative caching is disabled for in kube-dns](https://github.com/kubernetes/dns/issues/121).

Dnsmasq not caching negative responses may not seem like a big deal in a micro-services environment where processes are mostly querying for things that exist. But there is a quirk in Kubernetes pod DNS policy that results in a very high proportion of negative DNS responses.  The DNS policy causes a query for something like `google.com` to result in 3+ nonexistent name queries before actually trying `google.com` (e.g. `google.com.default.svc.cluster.local`, `google.com.svc.cluster.local`, `google.com.cluster.local`, `google.com.mylocal.com`). So in kube-dns, a query for `google.com` results in a series of queries of which only about 20% are cached, the other 80% have to make a comparatively expensive round trip to the kube-dns process.  In CoreDNS, all of the responses can be cached.  This can make a big difference in a cluster where many pods are communicating with one service.

## Test Environment

In a nutshell, I performed these tests using Kubernetes' DNS perf-test tool on a three node bare metal cluster using a private upstream DNS server. In more detail...

**Authoritative Base Kube-dns and CoreDNS Manifests**:  I have based the configurations of kube-dns and CoreDNS on their [authoritative base templates](https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/dns).  Specifically, tests use the authoritative default configurations for kube-dns (dnsmasq command flags in the Deployment) and CoreDNS (Corefile in the Configmap). A single instance of each DNS server was deployed.

**Kubernetes Cluster**: The tests were performed in a three node kubernetes cluster, using the Calico network plugin.  The Nodes were hosted at bare metal provider [Packet](https://www.packet.net), using a t1.small (Intel Xeon E3 4 Cores @ 3.5 GHz) as the master, and m1.xlarge (Intel Xeon E5 24 Cores @ 2.2 GHz) for the nodes.

**Kubernetes DNS Perf-test Tool**: I used the [kubernetes dns perf-tests tool](https://github.com/kubernetes/perf-tests/tree/master/dns) modified to [allow tests on CoreDNS](https://github.com/kubernetes/perf-tests/pull/114).  The tool by default runs the DNS server under test on the master node and the client pod on separate node (node 1).  The tests are not performed against the cluster's active cluster DNS, instead they configure, spin up and tear down an instance for each set of tests.

**Upstream Server**: The upstream server was a CoreDNS instance using the `template` plugin, configured to reply to all queries with answer of '1.2.3.4'.  It was run in the cluster, on a separate node from the client and DNS server under test (node 2).  This local server responds much more quickly than a real world upstream server would.

**DNS Server Versions**: CoreDNS 1.1.3, kube-dns 1.14.10


## The Performance Tests

I tested three categories of queries:

1. **invalid**: A name for a service that does not exist in the cluster, which prompts an `NXDOMAIN` response from the DNS service. e.g. `invalid-service.default.svc.cluster.local`
2. **service**: A name for a service that exists in the cluster. e.g. `kubernetes.default.svc.cluster.local`
3. **upstream**: A name outside of the cluster, which is forwarded upstream by the DNS service. e.g. `coredns.io`


### Comparisons on a Level Playing Field

Dnsmasq is single threaded, so when running on a multi core system, it will use at most one core. CoreDNS is multi threaded, so it could use all cores.  This presents an unfair advantage to CoreDNS when comparing single instances on a multi core system.  Scaling kube-dns instances horizontally in the test to level the playing field would be difficult to tune evenly.  So instead, I ran tests with restricted CPU allotted per container, such that the equivalent of one core at maximum would be used by the dns servers during the tests. The thinking here is that if kube-dns is deployed, it would/should be scaled appropriately such that it is not constrained by the single threaded dnsmasq. Testing multithreaded CoreDNS against an under-scaled kube-dns would be an unfair comparison.

Since CoreDNS operates in a single container, restricting CPU is easy.  But kube-dns contains two containers, so ensuring both its containers operate within one core without starving one is not as straight forward.  To account for this, I ran kube-dns CPU restrictions in two separate runs, one with all CPU allocated to dnsmasq (for tests of dnsmasq only tasks e.g. cache hits and upstream forwarding), and a second with kube-dns and dnsmasq sharing the CPU at it's nominal ratio (observationally, 70% kube-dns, 30% dnsmasq).


## Client Perspective Through the Cluster-First DNS Policy

Recall that in Kubernetes, queries for out-of-cluster names (e.g. google.com) result in *at least* 3 futile in-cluster queries, followed by the correct upstream query.  If there are any local search domains configured on the host node, then there can be more than three futile queries - and the additional queries can be even more expensive since they would typically be routed to an external DNS.  For this blog, I'm assuming that there are no local search domains configured on the host node.

Factoring in the Cluster-First DNS policy into the restricted CPU performance results, we can project performance from the client pod perspective using the following formulas...

Actual Out Of Cluster Cache Miss Latency = `(3 * invalid latency) + upstream latency`

Actual Out Of Cluster Cache Miss QPS = `3 / ((1 / invalid QPS) + (1 / upstream QPS))`

Actual Out Of Cluster Cache Hit Latency = `4 * cache-hit latency`

Actual Out Of Cluster Cache Hit QPS = `3 / ((1 / invalid QPS) + (1 / cache-hit QPS))`

So, from a client perspective CoreDNS performs better than kube-dns for out of cluster queries for both cache hits and misses.

To a lesser degree, the same problem exists for certain in-cluster queries. For example, in-cluster queries of the form `service-name.namespace` result in one invalid attempt per query.  The number of invalid attempts depends on how many domains in the search path the client must attempt before using the "correct" one. The number of invalid attempts per query type are...

* `service-name`: 0
* `service-name.namespace`: 1
* `service-name.namespace.svc`: 2
* `service-name.namespace.svc.cluster.local`: 3 (plus more if there are host node search domains)


## Conclusion

The following chart show the relative performance of CoreDNS vs Kube-DNS running in equivalent environment. Absolute numbers are not reported here because they were tested in a low-scale CPU restricted environment.  In other words, due to the method of testing the raw numbers are low, and I don't want them being taken out of context.  Each stacked bar shows the ratio of Kube-DNS to CoreDNS for that category.  For example, for out of cluster cached queries, CoreDNS performed 2X better than Kube-DNS, whereas for in-cluster cached queries, Kube-DNS performed about 17X better than CoreDNS.  Note that in-cluster cached queries here assumes that the service lies within the first search path (the same namespace as the client).

![latency chart](/images/CoreDNSvsKubeDNSPerformanceRatios.png)

The chart above is based on QPS ratios.  A chart showing latency ratios would show the inverse of the above.

In short, kube-dns is better at cached hits on in-cluster services, and CoreDNS is better in all other categories.  Whether or not CoreDNS will perform better in your cluster depends on your clusters DNS's load, cache hit ratio, and in-cluster out-cluster query mix.


