+++
title = "loadbalance"
description = "*loadbalance* randomize the order of A, AAAA and MX records."
weight = 24
tags = [ "plugin", "loadbalance" ]
categories = [ "plugin" ]
date = "2019-08-14T08:11:42.936486"
+++

## Description

The *loadbalance* will act as a round-robin DNS loadbalancer by randomizing the order of A, AAAA,
and MX records in the answer.

See [Wikipedia](https://en.wikipedia.org/wiki/Round-robin_DNS) about the pros and cons on this
setup. It will take care to sort any CNAMEs before any address records, because some stub resolver
implementations (like glibc) are particular about that.

## Syntax

~~~
loadbalance [POLICY]
~~~

* **POLICY** is how to balance, the default, and only option, is "round_robin".

## Examples

Load balance replies coming back from Google Public DNS:

~~~ corefile
. {
    loadbalance round_robin
    forward . 8.8.8.8 8.8.4.4
}
~~~
