+++
date = "2017-10-15T09:47:12Z"
title = "Autopath"
tags = ["Autopath","Kubernetes"]
description = "Server side search path extension with *autopath*"
author = "chris"
+++

CoreDNS 011 introduced a new plugin: [*autopath*](/plugins/autopath). It transparently addresses
a problem relating to search paths in kubernetes and can be useful in other scenarios. In
a nutshell, *autopath* helps reduce query load for Kubernetes service discovery. It does this by
resolving queries on the server side, instead of waiting for the client to request a search for each
domain in the search path one at a time.

In this post, I'll frame the problem by covering how DNS search paths work, and the way Kubernetes
uses search paths and how *autopath* feature can help resolve those problems.

## Search Paths

DNS resolvers allow for name resolution relative to their local domain. Clients can also be
configured with a list of domains to define a search path instead of just the local domain.  When
resolving a name, the client prefixes the name to each domain in the search path and queries the DNS
server until a successful response is returned. If none of the queries are successful then a search
on the absolute name is attempted.

The following example shows a search for `apple`,

```
% host -ta -v apple
Trying "apple.infoblox.com"
Trying "apple"
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 25023
;; flags: qr rd ra; QUERY: 1, ANSWER: 0, AUTHORITY: 1, ADDITIONAL: 0

;; QUESTION SECTION:
;apple.				IN	A

;; AUTHORITY SECTION:
apple.			889	IN	SOA	a0.nic.apple. noc.afilias-nst.info. 1000002523 10800 3600 2764800 900
```

Note that the first search tried is for `apple` in the local domain `infoblox.com`. This produced no
results, so the client then tried just plain `apple`, which resulted in an answer verifying the
existence of the `apple.` vanity top level domain.

In Unix, short name resolution behavior is more or less configured in `/etc/resolv.conf` with two
options, `search` and `ndots` paraphrased here:

- `search` Specifies the path of domains to search, in the order listed.  If not specified, the
  search list defaults to be the local domain of the host.
- `ndots` Sets a threshold for the number of dots which must appear in a name before an initial
  absolute query will be made. By default, `ndots` is set to 1, which means that a query for a name
  that contains one or more dots, such as `coredns.io` will skip the search path, and just query the
  name absolutely i.e. `coredns.io.`.

See the [resolv.conf(5)](http://man7.org/linux/man-pages/man5/resolv.conf.5.html) man page for more
detailed descriptions.

## Name Resolution in Kubernetes

Kubernetes has controls the resolv.conf configuration of pods using two different DNS Policies:
*ClusterFirst*, and *Default*.

- *ClusterFirst* - Causes the pod to use a special cluster oriented search path, enabling short name resolution.
- *Default* - Causes the pod to inherit the resolv.conf from the node it’s running on.

*ClusterFirst* is the default policy (ironically). In general, with exception of some Kubernetes
 infrastructure (notably the DNS service itself), all pods use the *ClusterFirst* policy.

The *ClusterFirst* will do the following:
- Use the cluster DNS service as the nameserver (e.g. coredns, or kube-dns)
- Use a search path that steps "out" from the local pod's namespace.
- Don't use the search path for searches containing 5 or more dots

For example, a pod in the `default`

```
nameserver 10.0.0.10
search default.svc.cluster.local svc.cluster.local cluster.local
options ndots:5
```

This search path allows a pod to query a service by name, e.g. `service1`, and get the IP of the
service in the same namespace. If the pod needs to find the a service of the same name in
a different namespace, the query needs to include the namespace. For example, `service1.ns2` will
produce the IP of `service1` in namespace `ns2` based on the second domain in the search path.

In addition to the three base domains, the search path configured on the pod’s node are added
(limited to 3 more). So, search path can get quite deep, up to 6 domains long.

The ClusterFirst DNS Policy also sets ndots to 5. This means that virtually every query a pod makes
is run through the search path. A query name would need to contain 5 dots to “skip” the search path.

So, why 5 ndots? The reason for this high ndots setting is the due to the potentially high number of
dots in a local service’s name.  For example, consider an SRV query for
`_http._tcp.service.namespace.svc`.  This query can be resolved using the 3rd domain in the search
path, trying `_http._tcp.service.namespace.svc.cluster.local`.  If ndots was set to something less
than 5, then the query would be the absolute name `_http._tcp.service.namespace.svc.` which would
not produce an answer.

## High ndots

The combination of long search path and high ndots means that almost every query made is eligible to
be searched on a long list of domains before finding an answer.  This has the most impact with
queries for external resources.  Any search for an external host with less than 5 dots will be run
through the whole search list before trying the absolute name.

This problem can manifest itself as scaling problems, such as high latency of DNS responses, high
load on the DNS server, and network congestion.

You can experiment with *autopath* yourself (even without Kubernetes), consider the following
Corefile:

~~~ corfile
. {
    proxy . 8.8.8.8
    autopath . resolv.conf
    log
}
~~~

And this file `resolv.conf` with a search path.

~~~ txt
search svc.cluster.local cluster.local local
~~~

For some extra logging in CoreDNS I'm using [this
patch](https://gist.github.com/miekg/7c5d4176c82a49e7b9a8d2d18249e421). Starting CoreDNS and asking
for `google.com.svc.cluster.local` we get the following answer:

~~~ txt
% dig @localhost A google.com.svc.cluster.local +noall +answer                                                          ~

google.com.svc.cluster.local. 299 IN	CNAME	google.com.
google.com.		299	IN	A	216.58.206.110
~~~

*autopath* performed the search path queries on the server and returned the answer with a CNAME
connecting the final answer to the original question asked.

In the logs below note the three unsuccessful queries before the final successful result.

~~~ txt
2017/10/15 09:36:00 [INFO] Querying google.com.svc.cluster.local.
2017/10/15 09:36:00 [INFO] NXDOMAIN for google.com.svc.cluster.local., continuing search
2017/10/15 09:36:00 [INFO] Querying google.com.cluster.local.
2017/10/15 09:36:00 [INFO] NXDOMAIN for google.com.cluster.local., continuing search
2017/10/15 09:36:00 [INFO] Querying google.com.local.
2017/10/15 09:36:00 [INFO] NXDOMAIN for google.com.local., continuing search
2017/10/15 09:36:00 [INFO] Querying google.com.
::1 - [15/Oct/2017:09:36:00 +0100] "A IN google.com.svc.cluster.local. udp 58 false 4096" NOERROR qr,rd,ra 98 94.736986ms
~~~
