+++
title = "health"
description = "*health* enables a health check endpoint."
weight = 15
tags = [ "plugin", "health" ]
categories = [ "plugin" ]
date = "2018-03-13T14:42:31.863930"
+++

## Description

By enabling *health* any plugin that implements
[healt.Healther interface](https://godoc.org/github.com/coredns/coredns/plugin/health#Healther)
will be queried for it's health. The combined health is exported, by default, on port 8080/health .

## Syntax

~~~
health [ADDRESS]
~~~

Optionally takes an address; the default is `:8080`. The health path is fixed to `/health`. The
health endpoint returns a 200 response code and the word "OK" when CoreDNS is healthy. It returns
a 503. *health* periodically (1s) polls plugin that exports health information. If any of the
plugin signals that it is unhealthy, the server will go unhealthy too. Each plugin that supports
health checks has a section "Health" in their README.

More options can be set with this extended syntax:

~~~
health [ADDRESS] {
    lameduck DURATION
}
~~~

* Where `lameduck` will make the process unhealthy then *wait* for **DURATION** before the process
  shuts down.

If you have multiple Server Block and need to export health for each of the plugins, you must run
health endpoints on different ports:

~~~ corefile
com {
    whoami
    health :8080
}

net {
    erratic
    health :8081
}
~~~

## Plugins

Any plugin that implements the Healther interface will be used to report health.

## Metrics

If monitoring is enabled (via the *prometheus* directive) then the following metric is exported:

* `coredns_health_request_duration_seconds{}` - duration to process a /health query. As this should
  be a local operation it should be fast. A (large) increases in this duration indicates the
  CoreDNS process is having trouble keeping up with its query load.

## Examples

Run another health endpoint on http://localhost:8091.

~~~ corefile
. {
    health localhost:8091
}
~~~

Set a lameduck duration of 1 second:

~~~ corefile
. {
    health localhost:8092 {
        lameduck 1s
    }
}
~~~
