+++
date = "2018-05-31T06:21:00Z"
description = "Comparing Performance of CoreDNS and Kube-DNS in Kubernetes."
tags = ["Kubernetes", "Performance", "Kube-DNS", "CoreDNS"]
title = "Performance Testing CoreDNS and Kube-dns"
author = "chris"
+++

# Performance Testing CoreDNS and Kube-dns

CoreDNS is on track to replace kube-dns as the the default cluster DNS service in Kubernetes.  Since Kubernetes 1.9, CoreDNS has been available as an optional replacement, and now CoreDNS is proposed to replace kube-dns as the default cluster DNS service in 1.12. As part of this transition, we have extended the existing kube-dns oriented [performance testing tool](https://github.com/kubernetes/perf-tests/tree/master/dns) to also [work with CoreDNS](https://github.com/kubernetes/perf-tests/pull/114).  This blog covers a performance comparison of kube-dns and CoreDNS using the updated tool.


## CoreDNS and Kube-dns Differences

I won't go into a comprehensive list of differences between CoreDNS and kube-dns here because most of them are not related to performance. But I’ll point out a few key differences that do weigh in on performance.

###C vs Go

CoreDNS is a single process, written in Go. Kube-dns actually operates as a few separate processes, the two we care about are *dnsmasq* which handles caching and routing upstream queries, and the namesake *kube-dns* which handles in-cluster names.  The *kube-dns* process is written in Go, but *dnsmasq* is written in C.  As it relates to performance, dnsmasq’s C code should be faster than CoreDNS’s Go implementation at caching and forwarding requests upstream. 

###Multithreading

Dnsmasq is single threaded, so its CPU usage is restricted to one core. CoreDNS is multi-threaded, to it can scale higher in a single instance by consuming multiple cores.  In kube-dns, all incoming DNS requests are routed through dnsmasq.  This presents a bottleneck which can be solved with horizontal scaling (i.e. load balanced replicas of kube-dns).  This horizontal scaling results in a larger memory footprint. It also complicates apples to apples testing a bit.

###Maximum Concurrent Queries

In kube-dns deployments, dnsmasq defaults to allowing a maximum of 150 concurrent queries.  If this number is exceeded, then dnsmasq will drop queries.      This can be a problem in a large scale deployment where you might see DNS request spikes (e.g. hundreds of pods hitting the DNS at nearly exactly the same time).  This can be adjusted up for environments that see big spikes in concurrent DNS requests.  CoreDNS does not have a similar hard limit for the number of concurrent queries.  In theory, instead of dropping packets, CoreDNS should handle big spikes, but with a higher latency.

###Negative Caching

Another performance related difference is how the two handle negative caching. Negative caching is a [standard](https://tools.ietf.org/html/rfc2308) that enables dns servers to cache responses that have no records.  In a DNS response, each record in the response will have a TTL.  A DNS cache uses this TTL to determine how long to keep the record in cache.  Empty responses, and denial of existence responses don’t normally have any records in them, so there is no TTL.  The standard way to handle this is to include an SOA record in the response with a TTL, so downstream caches know how long to keep the response.  In Kubernetes default deployments, CoreDNS keeps a negative cache, whereas [negative caching is disabled for in kube-dns](https://github.com/kubernetes/dns/issues/121).

Dnsmasq not caching negative responses may not seem like a big deal in a micro-services environment where processes are mostly querying for things that exist. But there is a quirk in Kubernetes pod DNS policy that results in a very high proportion of negative DNS responses.  The DNS policy causes a query for something like `google.com` to result in 3+ nonexistent name queries before actually trying `google.com` (e.g. `google.com.default.svc.cluster.local`, `google.com.svc.cluster.local`, `google.com.cluster.local`, `google.com.mylocal.com`). So in kube-dns, a query for `google.com` results in a series of queries of which only about 20% are cached, the other 80% have to make a comparatively expensive round trip to the kube-dns process.  In CoreDNS, all of the responses can be cached.  This can make a big difference in a cluster where many pods are communicating with one service.

## Test Environment

In a nutshell, I performed these tests using Kubernetes' DNS perf-test tool on a three node bare metal cluster using a private upstream DNS server. In more detail...

**Authoritative Base Kube-dns and CoreDNS Manifests**:  I have based the configurations of kube-dns and CoreDNS on their [authoritative base templates](https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/dns).  Specifically, tests use the authoritative default configurations for kube-dns (dnsmasq command flags in the Deployment) and CoreDNS (Corefile in the Configmap). A single instance of each DNS server was deployed.

**Kubernetes Cluster**: The tests were performed in a three node kubernetes cluster, using the Calico network plugin.  The Nodes were hosted at bare metal provider [Packet](https://www.packet.net), using a TBD as the master, and TBD for the nodes.

**Kubernetes DNS Perf-test Tool**: I used the [kubernetes dns perf-tests tool](https://github.com/kubernetes/perf-tests/tree/master/dns) modified to [allow tests on CoreDNS](https://github.com/kubernetes/perf-tests/pull/114).  The tool by default runs the DNS server under test on the master node and the client pod on separate node (node 1).  The tests are not performed against the cluster's active cluster DNS, instead they configure, spin up and tear down an instance for each set of tests.

**Upstream Server**: The upstream server was a CoreDNS instance using the `template` plugin, configured to reply to all queries with answer of '1.2.3.4'.  It was run in the cluster, on a separate node from the client and DNS server under test (node 2).  This local server responds much more quickly than a real world upstream server would.

**DNS Server Versions**: CoreDNS 1.1.3, kube-dns 1.14.10


## CoreDNS and Kube-dns QPS and Latency

Each set of tests is for 3 types of `A` record queries:

1. **invalid**: A name for a service that does not exist in the cluster, which prompts an `NXDOMAIN` response from the DNS service. e.g. `invalid-service.default.svc.cluster.local`
2. **service**: A name for a service that exists in the cluster. e.g. `kubernetes.default.svc.cluster.local`
3. **upstream**: A name outside of the cluster, which is forwarded upstream by the DNS service. e.g. `coredns.io`


### Results

TODO

## Cluster-First DNS Policy

Recall that in Kubernetes, queries for out-of-cluster names (e.g. google.com) result in at least 3 futile in-cluster queries, followed by the correct upstream query.

Factoring in the Cluster-First DNS policy into the restricted CPU performance results, we can project performance from the client pod perspective ...

### Out Of Cluster Cache Miss

TODO

### Out Of Cluster Cache Hit

In the cache hit scenario, for kube-dns we have to use the uncached metrics for negative queries since negative caching is disabled.

TODO


So, from a client perspective CoreDNS performs better than kube-dns for out of cluster queries for both cache hits and misses.  Kube-dns performance suffers greatly here because it disables negative caching.

## Conclusion

TODO


