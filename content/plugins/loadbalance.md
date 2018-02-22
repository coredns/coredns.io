+++
title = "loadbalance"
description = "*loadbalance* acts as a round-robin DNS loadbalancer by randomizing the order of A and AAAA records"
weight = 18
tags = [ "plugin", "loadbalance" ]
categories = [ "plugin" ]
date = "2018-02-22T08:55:16.404488"
+++
 in the answer.

## Description
 
 See [Wikipedia](https://en.wikipedia.org/wiki/Round-robin_DNS) about the pros and cons on this
 setup. It will take care to sort any CNAMEs before any address records, because some stub resolver
 implementations (like glibc) are particular about that.

## Syntax

~~~
loadbalance [POLICY]
~~~

* **POLICY** is how to balance, the default is "round_robin"

## Examples

Load balance replies coming back from Google Public DNS:

~~~ corefile
. {
    loadbalance round_robin
    proxy . 8.8.8.8 8.8.4.4
}
~~~
