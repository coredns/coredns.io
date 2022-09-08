+++
title = "health"
description = "*health* enables a health check endpoint."
weight = 24
tags = ["plugin", "health"]
categories = ["plugin"]
date = "2022-09-08T18:42:54.8775489"
+++

## Description

Enabled process wide health endpoint. When CoreDNS is up and running this returns a 200 OK HTTP
status code. The health is exported, by default, on port 8080/health.

## Syntax

~~~
health [ADDRESS]
~~~

Optionally takes an address; the default is `:8080`. The health path is fixed to `/health`. The
health endpoint returns a 200 response code and the word "OK" when this server is healthy.

An extra option can be set with this extended syntax:

~~~
health [ADDRESS] {
    lameduck DURATION
}
~~~

* Where `lameduck` will delay shutdown for **DURATION**. /health will still answer 200 OK.
  Note: The *ready* plugin will not answer OK while CoreDNS is in lame duck mode prior to shutdown.

If you have multiple Server Blocks, *health* can only be enabled in one of them (as it is process
wide). If you really need multiple endpoints, you must run health endpoints on different ports:

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

Doing this is supported but both endpoints ":8080" and ":8081" will export the exact same health.

## Metrics

If monitoring is enabled (via the *prometheus* plugin) then the following metrics are exported:

 * `coredns_health_request_duration_seconds{}` - The *health* plugin performs a self health check
    once per second on the `/health` endpoint. This metric is the duration to process that request.
    As this is a local operation it should be fast. A (large) increase in this
    duration indicates the CoreDNS process is having trouble keeping up with its query load.
 * `coredns_health_request_failures_total{}` - The number of times the self health check failed.

Note that these metrics *do not* have a `server` label, because being overloaded is a symptom of
the running process, *not* a specific server.

## Examples

Run another health endpoint on http://localhost:8091.

~~~ corefile
. {
    health localhost:8091
}
~~~

Set a lame duck duration of 1 second:

~~~ corefile
. {
    health localhost:8092 {
        lameduck 1s
    }
}
~~~
