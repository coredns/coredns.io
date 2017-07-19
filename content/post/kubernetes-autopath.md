+++
date = "2017-07-17T18:12:12Z"
title = "CoreDNS Kubernetes Autopath"
tags = ["CoreDNS","Kubernetes"]
slug = "CoreDNS Kubernetes Autopath"
author = "chris"
+++

CoreDNS 009 introduces a new feature in the Kubernetes middleware: Autopath. It transparently addresses a problem relating to search paths in kubernetes deployments. In a nutshell, Autopath helps reduce query load for Kubernetes service discovery. It does this by resolving short name queries on the server side, instead of waiting for the client to request a search for each domain in the search path one at a time.

In this blog post, I'll cover the background of how DNS search paths work, the problematic way Kubernetes uses search paths, and how the new Autopath feature helps to resolve those problems.

## Short name resolution and search paths

DNS resolvers allow for name resolution relative to their local domain. Clients can also be configured to search a path of domains instead of just the local domain.  When resolving a name, the client prefixes the name to each domain in the path and queries the DNS server until a successful response is returned. If none of the queries are successful, then a search on the absolute name is attempted.  

The following example shows a search for `apple`.

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

In Unix, short name resolution behavior is more or less configured in `/etc/resolv.conf` with two options, `search` and `ndots`

 - `search` Specifies the path of domains to search, in the order listed.  If not specified, the search list defaults to be the local domain of the host.
 - `ndots` Sets a threshold for the number of dots which must appear in a name before an initial absolute query will be made. By default, `ndots` is set to 1, which means that a query for a name that contains one or more dots, such as `coredns.io` will skip the search path, and just query the name absolutely i.e. `coredns.io.`.


## Short name resolution in Kubernetes

Kubernetes has controls the resolv.conf configuration of pods using two different DNS Policies: *ClusterFirst*, and *Default*.  

 - *ClusterFirst* - Causes the pod to use a special cluster oriented search path, enabling short name resolution.
 - *Default* - Causes the pod to inherit the resolv.conf from the node it’s running on.

*ClusterFirst* is the default policy (ironically). In general, with exception of some Kubernetes infrastructure (notably the DNS service itself), all pods use the *ClusterFirst* policy. It's the *ClusterFirst* that causes problems, so we'll focus on that policy here.

*ClusterFirst* creates a search path that steps "out" from the local pod's namespace. For example, for pods in the namespace `default` the search path would be:

```
default.svc.cluster.local  svc.cluster.local  cluster.local
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

Note the 3 unsuccessful queries before the final successful result. All of these attempts are complete round trip DNS queries, which puts more load on the system.  From the server's perspective, this looks like 3 independent queries from the same remote client.
With negative cache enabled, this also increases the number of cache entries in the DNS server, filling the cache with a significant percentage of nonsense entries that normally would not be searched for.  

This problem can manifest itself as scaling problems, such as high latency of DNS responses, high load on the DNS server, and network congestion.


## CoreDNS Autopath Solution 

The most direct way to resolve these issues would be reduce the number of queries required to arrive at the final answer, and this is precisely what the Autopath feature does.  When the Kubernetes middleware receives a query that looks like the first domain in the search path, Autopath anticipates the next questions and checks the remaining domains on the client's behalf until it finds an answer.  The client then receives an answer after making a single query, instead of 2-7 queries.  

Since the search path of a pod is predictable based on the DNS Policy and the namespace of the pod, Autopath can detect when a pod is using the first domain in the path and then continue searching the remaining domains in the path until a result if found.

The following diagram illustrates the basic flow of how the Autopath processes DNS queries.

![Kubernetes Autopath Flow](/images/autopath-flow.png)

**Diagram Notes**:  Highlighting three noteworthy aspects ...

1. The Autopath process only starts when a search would otherwise fail with a name error. If the search is successful, then there is no need to engage the Autopath function.  This means that for queries with an answer in the first search path, there is no per-query overhead.

2. Since the answer may be to a question different from the question the client originally asked, Autopath injects an artificial CNAME record to bridge the gap. Without this "magic glue" CNAME record, some clients will detect the mismatch between the answer and the question as a compromised response and throw an error.

3. Using a bit more magic, if none of the searches in the path produce an answer, and Autopath were to return an `NXDOMAIN` to the client, then the client would continue on with its own search path execution, cancelling the benefits of the Autopath feature. To prevent this from happening in this case, Autopath returns a successful result with no data, instead of a name error. Returning a success prompts the client to stop searching.


To illustrate with an example, here is the result of a query for `google.com` from a pod in namespace `default`.
   
```
# host -v -t a google.com 
Trying "google.com.default.svc.cluster.local"
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 50776
;; flags: qr aa rd; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;google.com.default.svc.cluster.local. IN A

;; ANSWER SECTION:
google.com.default.svc.cluster.local. 0	IN CNAME google.com.
google.com.		0	IN	A	192.0.2.53
```
   
The verbose output lets you see that the original query was `google.com`, but that the client first tries `google.com.default.svc.cluster.local`.  This is the short name `google.com` expanded with the first domain in the search path `default.svc.cluster.local`.  In the answer, there are two records.  The second record is the `A` record containing the IP address for `google.com.`. The first record is the `CNAME` "magic glue" which joins together the question section query name `google.com.default.svc.cluster.local.` with the `A` record.

On the server side, kubernetes middleware and the Autopath feature looked up the following on the backend before returning a response to the client...

|Search Domain|Question|Result|
|---|---|---|
|1st base domain|google.com.default.svc.cluster.local| fail |
|2nd base domain|google.com.svc.cluster.local| fail |
|3rd base domain|google.com.cluster.local| fail |
|host domain|google.com.coredns.io| fail |
|absolute|google.com.| success |

### Expanded vs fully qualified queries

A limitation of the feature relates to how Autopath detects the first search path.  Autopath  always assumes that if the first domain in the search path is present in a query, that must mean that the query was expanded by the client.  For example, consider the query `coredns.io.ns.svc.cluster.local.` received from a pod in the namespace `ns` Since the query is in `ns.svc.cluster.local.`, it is assumed that the original name queried was `coredns.io`, and that it was expanded with the first domain in the search path.  However, it's possible that the client queried for the fully qualified name `coredns.io.ns.svc.cluster.local.`.  From the perspective of the DNS server, these queries are identical.  Yet in the first case the client is following the search path, and in the second case the client is not.  Ideally, Autopath would search the path in the first case, and not in the second case.  But since Autopath cannot tell the difference, it defaults to assuming the first case.  In other words, it assumes that the original search was for the name `coredns.io` and expanded by the first search path, so therefore it should search the remaining domains in the path.

There is an optional parameter to autopath called NDOTS that allows you to tweak how the assumption is made.  If the number of dots in the prefix of the query if greater than or equal to the NDOTS value, then Autopath assumes the search was expanded on the client. If the number of dots in the prefix of the query is less than NDOTS, then its assumed that the query was fully qualified and not expanded. By default, NDOTS is 0.


## Configuring Autopath 

Configuring Kubernetes middleware to use Autopath is simple: just add the `autopath` keyword to the kubernetes section of your Corefile.  For example:

```
kubernetes cluster.local {
	autopath
}

``` 

While the default behavior will work for almost all cases, there are some optional settings you can use to adjust the default behavior.

```
kubernetes cluster.local {
	autopath [NDOTS[ RESPONSE[ RESOLV-CONF]]]
}
```

 - `NDOTS`: This sets a threshold of the number of dots required to use the Autopath feature. If the number of dots in the search prefix is at least NDOTS, then the Autopath feature will continue searching the prefix in each domain in the path.  Defaulting to zero, the optimization kicks in for all queries.  Example: If the Kubernetes middleware receives the query `foo.namespace.svc.cluster.local.`, the prefix `foo` is assumed to the original name, and the suffix `namespace.svc.cluster.local.` is assumed to be the first domain in the search path. Since `foo` contains 0 dots, and the default NDOTS value is 0, the remaining domains are searched.  If NDOTS was set to 1, then `foo.namespace.svc.cluster.local.` would not trigger the autopath feature, since `foo` contains less than one dot.
 
 - `RESPONSE`: This allows you to disable the optimization which prevents the client from continuing the search path if the absolute search name produces a name error.  Specifying `NXDOMAIN`, effectively disables the optimization. Specifying `SERVFAIL` will cause the feature to response with a server failure instead of a success.  `NOERROR` is the default value.  
 
 - `RESOLV-CONF`: If deploying CoreDNS outside of the cluster, then you may want to specify an alternate resolv.conf file to control which search domains are added to the base 3 domains.


