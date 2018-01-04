+++
title = "health"
description = "*health* enables a health check endpoint."
weight = 14
tags = [ "plugin", "health" ]
categories = [ "plugin" ]
date = "2018-01-04T12:51:22.157263"
+++

## Description

By enabling *health* any plugin that implements it will be queried for it's health. The combined
health is exported, by default, on port 8080/health .

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

Any plugin that implements the Healther interface will be used to report health.

## Examples

Run another health endpoint on http://localhost:8091.

~~~ corefile
. {
    health localhost:8091
}
~~~
