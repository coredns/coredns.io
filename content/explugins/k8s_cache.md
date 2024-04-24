+++
title = "k8s_cache"
description = "*k8s_cache* is a caching plugin with early refreshes for specified pods"
weight = 10
tags = [  "plugin" , "k8s", "cache" ]
categories = [ "plugin", "external" ]
date = "2024-04-24T15:20:00+02:00"
repo = "https://github.com/delta10/k8s_cache"
home = "https://github.com/delta10/k8s_cache#readme"
+++

## Description

This is a fork of [cache](https://github.com/coredns/coredns/tree/master/plugin/cache). It
adds an option to send a refreshed positive cache item first to pods with the label
`k8s-cache.coredns.io/early-refresh=true`. Other pods get it only after a specified
duration. This makes it possible to implement stable NetworkPolicy whitelists on the basis
of domain names that are resolved with DNS, using [Stable FQDNNetworkPolicies](https://github.com/delta10/fqdnnetworkpolicies).

The implementation uses an additional cache store called the "late cache", which is
shifted a number of seconds. On expiration, items in the late cache are replaced with
items from the early cache if they exist. When a request comes in, the plugin normally
checks first if the response is cached in the late cache, then in the early cache. If the
source IP matches a pod with the label `k8s-cache.coredns.io/early-refresh=true`, the late
cache is skipped and the early cache consulted immediately.

This plugin is intended as a replacement of the *cache* plugin and should not be used in
combination with it.

We will keep the code of this plugin in sync with *cache* as best as we can.

## Syntax

~~~ txt
k8s_cache [TTL] [ZONES...] {
    earlyrefresh [DURATION]
    success CAPACITY [TTL] [MINTTL]
    denial CAPACITY [TTL] [MINTTL]
    prefetch AMOUNT [[DURATION] [PERCENTAGE%]]
    serve_stale [DURATION] [REFRESH_MODE]
    servfail DURATION
    disable success|denial [ZONES...]
    keepttl
}
~~~

For details, see the [cache documentation](https://coredns.io/plugins/cache/). This plugin
adds one argument and changes the meaning of some other arguments slightly.

* `earlyrefresh` Set the **DURATION** (e.g., "5s") before which `early-refresh` pods get a
fresh reply. This option actually ***increases*** the cache duration of successful
responses for pods not having the early refresh label. Each client receives the current
cache duration *for it* as TTL response.
* `prefetch` Works as in *cache*, but it uses the expiration time of the early cache to
calculate whether prefetches should be done.
* `serve_stale` Works as in *cache*, but **DURATION** is counted from the expiration of
the early cache. For positive responses cached in the late cache, `serve_stale` starts
taking effect only when the late cache expires. After the late cache has expired, stale
serving will continue for **DURATION** minus the duration of `earlyrefresh`. Pods having
the early refresh label will never be served stale responses.

## Examples

Keep a positive and negative cache size of 10000 (default) and send cache refreshes 5
seconds earlier to pods with the early refresh label.

~~~ corefile
.:5300 {
  k8s_cache {
    success 10000
    denial 10000
    earlyrefresh 5s
  }
  forward . 8.8.8.8
}
~~~

For general caching examples, see the [cache documentation](https://coredns.io/plugins/cache/).
