+++
Categories = []
Description = ""
Keywords = []
Tags = []
date = "2016-08-06T21:00:36-07:00"
title = "Middlewares"
weight = 3
img = "chain.jpeg"
+++

Currently CoreDNS supports (among others) the following
[middlewares](https://github.com/coredns/coredns/tree/master/middleware):

* [chaos](https://github.com/coredns/coredns/tree/master/middleware/chaos/README.md): respond to CH
class queries
* [dnssec](https://github.com/coredns/coredns/tree/master/middleware/dnssec/README.md): on-the-fly
DNSSEC signing of records
* [etcd](https://github.com/coredns/coredns/tree/master/middleware/etcd/README.md): SkyDNS replacement
* [file](https://github.com/coredns/coredns/tree/master/middleware/file/README.md): serve DNS from a set
of files
* [health](https://github.com/coredns/coredns/tree/master/middleware/health/README.md): simple health
check
* [kubernetes](https://github.com/coredns/coredns/tree/master/middleware/kubernetes/README.md): use
CoreDNS as a KubeDNS replacement
* [loadbalance](https://github.com/coredns/coredns/tree/master/middleware/loadbalance/README.md):
shuffle A and AAAA records
* [metrics](https://github.com/coredns/coredns/tree/master/middleware/metrics/README.md):
[Prometheus](https://prometheus.io) metrics
* [pprof](https://github.com/coredns/coredns/tree/master/middleware/pprof/README.md): Go profiling
* [proxy](https://github.com/coredns/coredns/tree/master/middleware/proxy/README.md): forward queries to
an upstream (recursive) server
* [rewrite](https://github.com/coredns/coredns/tree/master/middleware/rewrite/README.md): rewrite
incoming queries
* [secondary](https://github.com/coredns/coredns/tree/master/middleware/secondary/README.md): be
a secondary nameserver and retrieve zones from a primary
