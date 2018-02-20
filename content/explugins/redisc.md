+++
title = "redisc"
description = "*redisc* - enables a networked cache using Redis."
weight = 10
tags = [  "plugin" , "redisc" ]
categories = [ "plugin", "external" ]
date = "2018-02-17T19:36:00+00:00"
repo = "https://github.com/miek/redis"
home = "https://github.com/miek/redis/blob/master/README.md"
+++

## Description

With *redisc* responses can be cached for up to 3600s. Caching in Redis is mostly usefull in
a setup where multiple CoreDNS instances share a VIP. E.g. multiple CoreDNS pods in a Kubernetes
cluster.

If Redis is not reacheable this plugin will be a noop. The *cache* and *redisc* plugin can be used
together, where *cache* is the L1 and *redisc* is the L2 level cache.
If multiple CoreDNS instances get a cache miss for the same item, they will all be fetching the same
information from an upstream and updating the cache, i.e. there is no (extra) coordination between
those instances.

## Syntax

~~~ txt
redisc [TTL] [ZONES...]
~~~

* **TTL** max TTL in seconds. If not specified, the maximum TTL will be used, which is 3600 for
    noerror responses and 1800 for denial of existence ones.
    Setting a TTL of 300: `redisc 300` would cache records up to 300 seconds.
* **ZONES** zones it should cache for. If empty, the zones from the configuration block are used.

Each element in the Redis cache is cached according to its TTL (with **TTL** as the max). For the negative
cache, the SOA's MinTTL value is used. As no endpoint is specfied the default of `127.0.0.1:6379` will
be used.

If you want more control:

~~~ txt
redisc [TTL] [ZONES...] {
    endpoint ENDPOINT
    success TTL
    denial TTL
}
~~~

* **TTL**  and **ZONES** as above.
* `endpoint` specifies which **ENDPOINT** to use for Redis, this default to `127.0.0.1:6379`.
* `success`, override the settings for caching successful responses. **TTL** overrides the cache maximum TTL.
* `denial`, override the settings for caching denial of existence responses. **TTL** overrides the cache maximum TTL.
  There is a third category (`error`) but those responses are never cached.

## Metrics

If monitoring is enabled (via the *prometheus* directive) then the following metrics are exported:

* `coredns_redisc_hits_total{}` - Counter of cache hits.
* `coredns_redisc_misses_total{}` - Counter of cache misses.

## Examples

Enable caching for all zones, cache locally and also cache for up to 40s in the cluster wide Redis.

~~~ corefile
. {
    cache 30
    redisc 40 {
        endpoint 10.0.240.1:69
    }
    whoami
}
~~~

Proxy to Google Public DNS and only cache responses for example.org (and below).

~~~ corefile
. {
    proxy . 8.8.8.8:53
    redisc example.org
}
~~~

## See Also

See [the Redis site for more information](https://redis.io) on Redis. An external plugin called
[redis](https://coredns.io/explugins/redis) already exists, hence this is named *redisc*, for
"redis cache".

## Bugs

There is little unit testing.
