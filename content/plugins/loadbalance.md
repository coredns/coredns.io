+++
title = "loadbalance"
description = "*loadbalance* randomizes the order of A, AAAA and MX records."
weight = 26
tags = ["plugin", "loadbalance"]
categories = ["plugin"]
date = "2020-02-06T11:07:03.877382"
+++

## Description

The *loadbalance* will act as a round-robin DNS load balancer by randomizing the order of A, AAAA,
and MX records in the answer.

See [Wikipedia](https://en.wikipedia.org/wiki/Round-robin_DNS) about the pros and cons of this
setup. It will take care to sort any CNAMEs before any address records, because some stub resolver
implementations (like glibc) are particular about that.

## Syntax

~~~
loadbalance [POLICY]
~~~

* **POLICY** is how to balance. The default, and only option, is "round_robin".

## Examples

Load balance replies coming back from Google Public DNS:

~~~ corefile
. {
    loadbalance round_robin
    forward . 8.8.8.8 8.8.4.4
}
~~~
