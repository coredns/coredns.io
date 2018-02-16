# What is CoreDNS

CoreDNS is a DNS server. It is written in Go. What you say, another DNS server? Why should it exist,
what makes CoreDNS special?

CoreDNS is different from other DNS servers, such as the (all excellent) BIND, Knot, PowerDNS and
Unbound (technically a resolver, but still worthy a mention), because by default it can't do
anything. CoreDNS is powered by plugins. Plugins make CoreDNS a DNS server.

Plugins can be stand alone or work together to perform a "DNS function". The are currently about 30
plugins included in the default CoreDNS install.

What's a "DNS function" you ask? For the purpose of CoreDNS we define it as a piece of software that
generates a response to the query. The form these take can wildly deviate, there are plugins that
only create (Prometheus) metrics or cache a result and thus don't create a response by them selves.
There are plugin that communicate to a Kubernetes backend and provide service discovery inside
a Kubernetes cluster. There a plugins that read data from a file or database.

## Current Focus

* Kubernetes
* Rock solid forwarding
* Authoritative DNS

## Origin

Caddy, bladdy, included some older blog posts in a condensed form.
