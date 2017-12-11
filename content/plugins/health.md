+++
title = "health"
description = "*health* enables a simple health check endpoint."
weight = 14
tags = [ "plugin", "health" ]
categories = [ "plugin" ]
date = "2017-12-11T16:50:50.553018"
+++

By default, it listens on port 8080.

## Syntax

~~~
health [ADDRESS]
~~~

Optionally takes an address; the default is `:8080`. The health path is fixed to `/health`. The
health endpoint returns a 200 response code and the word "OK" when CoreDNS is healthy. It returns
a 503. *health* periodically (1s) polls plugin that exports health information. If any of the
plugin signals that it is unhealthy, the server will go unhealthy too. Each plugin that
supports health checks has a section "Health" in their README.

## Plugins

The following plugins report health to the health plugin:

* erratic
* kubernetes

## Examples

Run another health endpoint on http://localhost:8091.

~~~ corefile
. {
    health localhost:8091
}
~~~
