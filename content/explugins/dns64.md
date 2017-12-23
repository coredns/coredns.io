+++
title = "dns64"
description = "*dns64* implements the DNS64 IPv6 transition mechanism."
weight = 10
tags = [  "plugin" , "dns64" ]
categories = [ "plugin", "external" ]
date = "2017-12-24T05:17:00+08:00"
repo = "https://github.com/serverwentdown/dns64"
home = "https://github.com/serverwentdown/dns64/blob/master/README.md"
+++

The *dns64* plugin implements the DNS64 IPv6 transition mechanism. From Wikipedia:

> DNS64 describes a DNS server that when asked for a domain's AAAA records, but only finds
> A records, synthesizes the AAAA records from the A records.

The synthesis in only performed if the query came in via IPv6.

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

~~~
dns64 {
    upstream 8.8.8.8:53
}
proxy . 8.8.8.8:53
~~~

Perform dns64 AAAA synthesizing using 8.8.8.8 for resolving any A records.

## See Also

<https://en.wikipedia.org/wiki/IPv6_transition_mechanism#DNS64> and RFC 6147.
