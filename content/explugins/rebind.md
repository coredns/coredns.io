+++
title = "rebind"
description = "*rebind* - rebinds domain names to a different IP address."
weight = 10
tags = [  "plugin" , "rebind" ]
categories = [ "plugin", "external" ]
date = "2024-06-17T00:00:00+00:00"
repo = "https://github.com/ivantsepp/coredns-rebind"
home = "https://github.com/ivantsepp/coredns-rebind/blob/main/README.md"
+++

## Description

The rebind plugin rebinds a domain from one IP address to another IP address to facilitate testing [DNS Rebinding vulnerabilities](https://en.wikipedia.org/wiki/DNS_rebinding).

## Syntax

~~~ corefile
rebind example.com {
  first_ip 1.2.3.4
  second_ip 0.0.0.0
  strategy first_then_second
}
~~~

- **first_ip** is the first IP address. This is usually an IP address that you own
- **second_ip** is the second IP address to rebind to. This is usually the target IP address of the vulnerable server
- **strategy** is one of the following:
  - first_then_second: responds with the `first_ip` and then responds with the `second_ip` address for all subsequent requests
  - random: responds with a random selection of `first_ip` and `second_ip`
  - round_robin: responds in a round robin fashion of `first_ip` and then `second_ip`

## Examples

In this configuration, a DNS request to `rebind.example.com` will receive a response of `1.2.3.4`. All future DNS requests will respond with `0.0.0.0`.

~~~ corefile
example.com {
  rebind rebind.example.com {
    first_ip 1.2.3.4
    second_ip 0.0.0.0
  }
}
~~~

## Also See

See the [manual](https://coredns.io/manual).
