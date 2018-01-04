+++
title = "pprof"
description = "*pprof* publishes runtime profiling data at endpoints under `/debug/pprof`."
weight = 21
tags = [ "plugin", "pprof" ]
categories = [ "plugin" ]
date = "2018-01-04T12:51:22.160691"
+++

## Description

You can visit `/debug/pprof` on your site for an index of the available endpoints. By default it
will listen on localhost:6053.

> This is a debugging tool. Certain requests (such as collecting execution traces) can be slow. If
> you use pprof on a live server, consider restricting access or enabling it only temporarily.

For more information, please see [Go's pprof
documentation](https://golang.org/pkg/net/http/pprof/) and read
[Profiling Go Programs](https://blog.golang.org/profiling-go-programs).

## Syntax

~~~
pprof [ADDRESS]
~~~

If not specified, ADDRESS defaults to localhost:6053.

## Examples

Enable pprof endpoints:

~~~
. {
    pprof
}
~~~

Listen on an alternate address:

~~~ txt
. {
    pprof 10.9.8.7:6060
}
~~~

Listen on an all addresses on port 6060:

~~~ txt
. {
    pprof :6060
}
~~~

# Also See

See [Go's pprof documentation](https://golang.org/pkg/net/http/pprof/) and [Profiling Go
Programs](https://blog.golang.org/profiling-go-programs).
