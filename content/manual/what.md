# What is CoreDNS

CoreDNS is a DNS server. It is written in [Go](https://golang.org).

CoreDNS is different from other DNS servers, such as the (all excellent) BIND, Knot, PowerDNS and
Unbound (technically a resolver, but still worthy a mention), because it is very flexible; as it
chains plugins. CoreDNS is powered by plugins.

Plugins can be stand alone or work together to perform a "DNS function".

What's a "DNS function" you ask? For the purpose of CoreDNS we define it as a piece of software that
implements the CoreDNS Plugin API. The form these take can wildly deviate, there are plugins that
create (Prometheus) metrics or cache a result and thus don't create a response by them selves.
There are plugins that communicate with a Kubernetes backend and provide service discovery inside
a Kubernetes cluster. Then there a plugins that read data from a file or database.

The are currently about 30 plugins included in the default CoreDNS install, and a whole bunch of
[external](https://coredns.io/explugins) plugins.

## Current Focus

* Documentation
* Kubernetes
* Rock solid forwarding
* Authoritative DNS

## Origin

Not sure; dealt in blog items.
