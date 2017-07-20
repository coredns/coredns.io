+++
date = "2017-07-17T18:12:12Z"
title = "CoreDNS Kubernetes Autopath pt 1"
tags = ["CoreDNS","Kubernetes"]
slug = "CoreDNS Kubernetes Autopath pt 1"
author = "chris"
+++

CoreDNS 010 introduces a new feature in the Kubernetes middleware: Autopath. It transparently addresses a problem relating to search paths in kubernetes deployments. In a nutshell, Autopath helps reduce query load for Kubernetes service discovery. It does this by resolving short name queries on the server side, instead of waiting for the client to request a search for each domain in the search path one at a time.

In this blog post, I'll frame the problem by covering how DNS search paths work, and the problematic way Kubernetes uses search paths.  In a following post, I'll explain how you can use the new Autopath feature to help resolve those problems.


## Short name resolution and search paths

DNS resolvers allow for name resolution relative to their local domain. Clients can also be configured with a list of domains to define a search path instead of just the local domain.  When resolving a name, the client prefixes the name to each domain in the search path and queries the DNS server until a successful response is returned. If none of the queries are successful then a search on the absolute name is attempted.  

The following example shows a search for `apple`, 

```
$ host -ta -v apple
Trying "apple.infoblox.com"
Trying "apple"
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 25023
;; flags: qr rd ra; QUERY: 1, ANSWER: 0, AUTHORITY: 1, ADDITIONAL: 0

;; QUESTION SECTION:
;apple.				IN	A

;; AUTHORITY SECTION:
apple.			889	IN	SOA	a0.nic.apple. noc.afilias-nst.info. 1000002523 10800 3600 2764800 900

Received 86 bytes from 10.102.3.10#53 in 91 ms
```

Note that the first search tried is for `apple` in the local domain `infoblox.com`. This produced no results, so the client then tried just plain `apple`, which resulted in an answer verifying the existence of the `apple.` vanity top level domain.

In Unix, short name resolution behavior is more or less configured in `/etc/resolv.conf` with two options, `search` and `ndots` paraphrased here ...

 - `search` Specifies the path of domains to search, in the order listed.  If not specified, the search list defaults to be the local domain of the host.
 - `ndots` Sets a threshold for the number of dots which must appear in a name before an initial absolute query will be made. By default, `ndots` is set to 1, which means that a query for a name that contains one or more dots, such as `coredns.io` will skip the search path, and just query the name absolutely i.e. `coredns.io.`.

See the [resolv.conf(5)](http://man7.org/linux/man-pages/man5/resolv.conf.5.html) man page for more detailed descriptions.

## Short name resolution in Kubernetes

Kubernetes has controls the resolv.conf configuration of pods using two different DNS Policies: *ClusterFirst*, and *Default*.  

 - *ClusterFirst* - Causes the pod to use a special cluster oriented search path, enabling short name resolution.
 - *Default* - Causes the pod to inherit the resolv.conf from the node it’s running on.

*ClusterFirst* is the default policy (ironically). In general, with exception of some Kubernetes infrastructure (notably the DNS service itself), all pods use the *ClusterFirst* policy.

The *ClusterFirst* will do the following:
- Use the cluster DNS service as the nameserver (e.g. coredns, or kube-dns)
- Use a search path that steps "out" from the local pod's namespace.
- Dont use the search path for searches contianing 5 or more dots

For example, a pod in the `default`

```
nameserver 10.0.0.10
search default.svc.cluster.local svc.cluster.local cluster.local
options ndots:5
```

This search path allows a pod to query a service by name, e.g. `service1`, and get the IP of the service in the same namespace. If the pod needs to find the a service of the same name in a different namespace, the query needs to include the namespace. For example, `service1.ns2` will produce the IP of `service1` in namespace `ns2` based on the second domain in the search path.

In addition to the three base domains, the search path configured on the pod’s node are added (limited to 3 more). So, search path can get quite deep, up to 6 domains long.

The ClusterFirst DNS Policy also sets ndots to 5. This means that virtually every query a pod makes is run through the search path. A query name would need to contain 5 dots to “skip” the search path. 

So, why 5 ndots? The reason for this high ndots setting is the due to the potentially high number of dots in a local service’s name.  For example, consider an SRV query for `_http._tcp.service.namespace.svc`.  This query can be resolved using the 3rd domain in the search path, trying `_http._tcp.service.namespace.svc.cluster.local`.  If ndots was set to something less than 5, then the query would be the absolute name `_http._tcp.service.namespace.svc.` which would not produce an answer.

The most dots possible in a Kubernetes short name is 6. A name like that would look something like `_http._tcp.endpoint.service.namespace.federation.svc`.  To permit this name to be checked in the search path, ndots would have to be set to 7. This is a `SRV` record query of an endpoint of a federated service.

If you omit queries for `SRV` records, and disregard federation, the most number of dots you would have in a local service would be 3 `endpoint.service.namespace.svc`.  So an ndots of 4 would probably suffice 99% of the time.  Presumably, an ndots of 5 was considered a fair middle ground. 

In summary, the Kubernetes DNS Policy configures a long search path and high ndots for pods.  In the next section I'll cover why this is a problematic combination.


## Long search paths and high ndots

The combination of long search path and high ndots means that almost every query made is eligible to be searched on a long list of domains before finding an answer.  This has the most impact with queries for external resources.  Any search for an external host with less than 5 dots will be run through the whole search list before trying the absolute name. For example, a query for `coredns.io` will result in the following ...

```
# host -v -ta coredns.io
Trying "coredns.io.default.svc.cluster.local"
Trying "coredns.io.svc.cluster.local"
Trying "coredns.io.cluster.local"
Trying "coredns.io"
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 8985
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;coredns.io.			IN	A

;; ANSWER SECTION:
coredns.io.		600	IN	A	176.58.119.54
```

Note the three unsuccessful queries before the final successful result. All of these attempts are complete round trip DNS queries, which puts more load on the system (both client an server).  From the server's perspective, this looks like three independent queries from the same remote client.
With negative cache enabled, this also increases the number of cache entries in the DNS server, filling the cache with a significant percentage of nonsense entries that normally would not be searched for.  

This problem can manifest itself as scaling problems, such as high latency of DNS responses, high load on the DNS server, and network congestion.