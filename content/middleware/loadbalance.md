+++
title = "loadbalance"
description = "*loadbalance* acts as a round-robin DNS loadbalancer by randomizing the order of A and AAAA records  in the answer.    See [Wikipedia](https://en.wikipedia.org/wiki/Round-robin_DNS) about the pros and cons on this  setup. It will take care to sort any CNAMEs before any address records, because some stub resolver  implementations (like glibc) are particular about that. "
weight = 14
tags = [  "middleware" , "loadbalance" ]
categories = [ "middleware" ]
date = "2017-07-24T15:25:40+00:00"
+++

## Syntax

~~~
loadbalance [POLICY]
~~~

* **POLICY** is how to balance, the default is "round_robin"

## Examples

~~~
loadbalance round_robin
~~~

