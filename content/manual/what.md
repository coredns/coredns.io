# What is CoreDNS

CoreDNS is a DNS server. It is written in [Go](https://golang.org).

CoreDNS is different from other DNS servers, such as the (all excellent)
[BIND](https://www.isc.org/blogs/category/bind/),
[Knot](https://www.knot-dns.cz/),
[PowerDNS](https://www.powerdns.com/) and
[Unbound](https://www.unbound.net/) (technically a resolver, but still worthy a mention), because it
is very flexible; it chains plugins.

Plugins can be stand alone or work together to perform a "DNS function".

So what's a "DNS function"? For the purpose of CoreDNS we define it as a piece of software that
implements the CoreDNS Plugin API. The functionally implemented can wildly deviate. There are
plugins that create only (Prometheus) [metrics](/plugins/metrics) or [cache](/plugin/cache) a result
and thus don't create a response directly. Then there are plugins that communicate with
a [Kubernetes backend](/plugins/kubernetes) and provide service discovery inside a Kubernetes
cluster. There a plugins that read data from a [file](/plugins/file) or [database](/explugins/pdsql).

The are currently about 30 plugins included in the default CoreDNS install, and a whole bunch of
[external](/explugins) plugins.

> CoreDNS is powered by plugins.

## Reloading process

TODO

## Metrics

TODO
