+++
title = "multicluster"
description = "*multicluster* plugin is an implementation of Multicluster DNS specification."
weight = 10
tags = [  "plugin" , "multicluster" ]
categories = [ "plugin", "external" ]
date = "2021-11-09T12:37:19+01:00"
repo = "https://github.com/coredns/multicluster/"
home = "https://github.com/coredns/multicluster#readme"
+++

## Description

The *multicluster* plugin implements the [Kubernetes DNS-Based Multicluster Service Discovery
Specification](https://github.com/kubernetes/enhancements/pull/2577).

## Syntax

```
multicluster [ZONES...] {
    kubeconfig KUBECONFIG [CONTEXT]
    noendpoints
    fallthrough [ZONES...]
}
```

* `kubeconfig` **KUBECONFIG [CONTEXT]** authenticates the connection to a remote k8s cluster using a kubeconfig file. **[CONTEXT]** is optional, if not set, then the current context specified in kubeconfig will be used. It supports TLS, username and password, or token-based authentication. This option is ignored if connecting in-cluster (i.e., the endpoint is not specified).
* `noendpoints` will turn off the serving of endpoint records by disabling the watch on endpoints. All endpoint queries and headless service queries will result in an NXDOMAIN.
* `fallthrough` **[ZONES...]** If a query for a record in the zones for which the plugin is authoritative results in NXDOMAIN, normally that is what the response will be. However, if you specify this option, the query will instead be passed on down the plugin chain, which can include another plugin to handle the query. If **[ZONES...]** is omitted, then fallthrough happens for all zones for which the plugin is authoritative. If specific zones are listed (for example `in-addr.arpa` and `ip6.arpa`), then only queries for those zones will be subject to fallthrough.

## Example

Handle all queries in the `clusterset.local` zone. 

```
.:53 {
    multicluster clusterset.local
}
```
