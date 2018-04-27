+++
title = "example"
description = "*example* - prints 'example' on every query received."
weight = 10
tags = [  "plugin" , "example" ]
categories = [ "plugin", "external" ]
date = "2017-07-25T21:57:00+08:00"
repo = "https://github.com/coredns/example"
home = "https://github.com/coredns/example/blob/master/README.md"
+++

## Description

The example plugin prints "example" on every query received. It serves as documentation for
writing CoreDNS plugins.

## Syntax

~~~ txt
example
~~~

## Metrics

If monitoring is enabled (via the *prometheus* directive) the following metric is exported:

* `coredns_example_request_count_total{server}` - query count to the *example* plugin.

The `server` label indicated which server handled the request, see the *metrics* plugin for details.

## Health

This plugin implements dynamic health checking. It will always return healthy though.

## Examples

In this configuration, we forward all queries to 9.9.9.9 and print "example" whenever we receive
a query.

``` corefile
. {
  forward . 9.9.9.9
  example
}
```

## Also See

See the [manual](https://coredns.io/manual).
