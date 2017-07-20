+++
date = "2017-07-18T18:12:12Z"
title = "CoreDNS Kubernetes Autopath pt 2"
tags = ["CoreDNS","Kubernetes"]
slug = "CoreDNS Kubernetes Autopath pt 2"
author = "chris"
+++

In the the [last post](/2017/07/17/coredns-kubernetes-autopath-pt-1/) we learned about how DNS search paths work, and the problematic way Kubernetes uses them.
In this post, I’ll explain how you can use the Autopath feature introduced in CoreDNS 010 to help resolve those problems.

## CoreDNS Autopath Solution 

The most direct way to resolve the Kubernetes search path problem would be reduce the number of queries required to arrive at the final answer. This is precisely what the Autopath feature does.  When the Kubernetes middleware receives a query that looks like the first domain in the search path, Autopath anticipates the next questions and checks the remaining domains on the client's behalf until it finds an answer.  The client then receives an answer after making a single query, instead of 2-7 queries.  

Since the search path of a pod is predictable based on the DNS Policy and the namespace of the pod, Autopath can detect when a pod is using the first domain in the path and then continue searching the remaining domains in the path until a result if found.

The following diagram illustrates the basic flow of how the Autopath processes DNS queries.

![Kubernetes Autopath Flow](/images/autopath-flow.png)

I want to highlight three noteworthy aspects in the design...

1. The Autopath process only starts when a search would otherwise fail with a name error. If the search is successful, then there is no need to engage the Autopath function.  This means that for queries with an answer in the first search path, there is no per-query overhead incurred by enabling the Autopath feature.

2. Since the answer a client receives may be to a question different from the question the client originally asked, Autopath injects an artificial CNAME record to bridge the gap. Without this "magic glue" CNAME record, some clients will detect the mismatch between the answer and the question as a compromised response and throw an error.

3. Using a bit more magic, if none of the searches in the path produce an answer, and Autopath were to return an `NXDOMAIN` to the client, then the client would continue on with its own search path execution, defeating the benefits of the Autopath feature. To prevent this from happening, Autopath returns a successful result with no data, instead of a name error. Returning a success prompts the client to stop searching.


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

There is an optional parameter to autopath called `ndots` that allows you to tweak how the assumption is made.  If the number of dots in the prefix of the query if greater than or equal to the `ndots` value, then Autopath assumes the search was expanded on the client. If the number of dots in the prefix of the query is less than `ndots`, then its assumed that the query was fully qualified and not expanded. By default, `ndots` is 0.


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
	autopath [ndots[ response[ resolve-conf]]]
}
```

 - `ndots`: This sets a threshold of the number of dots required in a query name to use the Autopath feature. If the number of dots in the search _prefix_ is at least NDOTS, then the Autopath feature will continue searching the prefix in each domain in the path.  Defaulting to zero, the optimization kicks in for all queries.  Example: If the Kubernetes middleware receives the query `foo.namespace.svc.cluster.local.`, the prefix `foo` is assumed to be the original name, and the suffix `namespace.svc.cluster.local.` is assumed to be the first domain in the search path. Since `foo` contains 0 dots, and the default NDOTS value is 0, the remaining domains are searched.  If `ndots` was set to 1, then `foo.namespace.svc.cluster.local.` would not trigger the autopath feature, since `foo` contains less than one dot.
 
 - `response`: This allows you to disable the optimization which prevents the client from continuing the search path if the absolute search name produces a name error.  Specifying `NXDOMAIN`, effectively disables the optimization. Specifying `SERVFAIL` will cause the feature to response with a server failure instead of a success.  `NOERROR` is the default value.  
 
 - `resolv-conf`: If deploying CoreDNS outside of the cluster, then you may want to specify an alternate resolv.conf file to control which search domains are added to the base 3 domains.

