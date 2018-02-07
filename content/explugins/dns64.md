+++
title = "dns64"
description = "*dns64* - implement the DNS64 IPv6 transition mechanism."
weight = 10
tags = [  "plugin" , "dns64" ]
categories = [ "plugin", "external" ]
date = "2017-12-24T05:17:00+08:00"
repo = "https://github.com/serverwentdown/dns64"
home = "https://github.com/serverwentdown/dns64/blob/master/README.md"
+++

## Description

The *dns64* plugin implements the DNS64 IPv6 transition mechanism. From Wikipedia:

> DNS64 describes a DNS server that when asked for a domain's AAAA records, but only finds
> A records, synthesizes the AAAA records from the A records.

The synthesis is only performed if the query came in via IPv6.

## TODO

Not all features required by DNS64 are implemented, only basic AAAA synthesis.

* [ ] Resolve PTR records
* [ ] Follow CNAME records
* [ ] Make resolver DNSSEC aware

## Syntax

~~~
dns64 {
    upstream ADDRESS...
    prefix IPV6
}
~~~

* `upstream` specifies the upstream resolver.
* `prefix` specifies any local IPv6 prefix to use, in addition to the well known
  prefix (64:ff9b::/96).

## Examples

In recursive resolver mode:

~~~
# Perform dns64 AAAA synthesizing using 8.8.8.8 for resolving any A
dns64 {
    upstream 8.8.8.8:53
}
proxy . 8.8.8.8:53
~~~

To make DNS64 resolve authoritatively, do:

~~~
dns64 {
    upstream localhost:53
    # caveat: additional round trip through networking stack
}
file example.com.db
~~~

## See Also

<https://en.wikipedia.org/wiki/IPv6_transition_mechanism#DNS64> and RFC 6147.
