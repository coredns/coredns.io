+++
title = "health"
description = "This module enables a simple health check endpoint. By default it will listen on port 8080. "
weight = 11
tags = [  "middleware" , "health" ]
categories = [ "middleware" ]
date = "2017-07-24T15:25:40+00:00"
+++

## Syntax

~~~
health [ADDRESS]
~~~

Optionally takes an address; the default is `:8080`. The health path is fixed to `/health`. It
will just return "OK" when CoreDNS is healthy, which currently mean: it is up and running.

This middleware only needs to be enabled once.

## Examples

~~~
health localhost:8091
~~~

