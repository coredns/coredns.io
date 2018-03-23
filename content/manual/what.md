# What is CoreDNS

CoreDNS is a DNS server. It is written in [Go](https://golang.org).

CoreDNS is different from other DNS servers, such as (all excellent)
[BIND](https://www.isc.org/blogs/category/bind/),
[Knot](https://www.knot-dns.cz/),
[PowerDNS](https://www.powerdns.com/) and
[Unbound](https://www.unbound.net/) (technically a resolver, but still worth a mention), because it
is very flexible; it chains plugins.

Plugins can be stand alone or work together to perform a "DNS function".

So what's a "DNS function"? For the purpose of CoreDNS we define it as a piece of software that
implements the CoreDNS Plugin API. The functionality implemented can wildly deviate. There are
plugins that don't themselves create a response, such as [metrics](/plugins/metrics) or
[cache](/plugin/cache a result), but add functionality. Then there are plugins that *do* generate
a response. These can also do anything; there are plugins that communicate with
[Kubernetes](/plugins/kubernetes) and provide service discovery, plugins that read data from
a [file](/plugins/file) or [database](/explugins/pdsql).

There are currently about 30 plugins included in the default CoreDNS install, but there are also a whole
bunch of [external](/explugins) plugins that you can compile into CoreDNS to extend its
functionality.

> CoreDNS is powered by plugins.

Writing new [plugins](#writing-plugins) should be easy enough, but requires knowing Go and having
some insight into how DNS works. CoreDNS abstracts away all other bits, so that you can just
focus on writing the plugin functionality you need.
