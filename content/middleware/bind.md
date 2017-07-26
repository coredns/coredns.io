+++
title = "bind"
description = "*bind* overrides the host to which the server should bind. "
weight = 2
tags = [  "middleware" , "bind" ]
categories = [ "middleware" ]
date = "2017-07-26T08:45:58+01:00"
+++

Normally, the listener binds to the wildcard host. However, you may force the listener to bind to
another IP instead. This directive accepts only an address, not a port.

## Syntax

~~~ txt
bind ADDRESS
~~~

**ADDRESS** is the IP address to bind to.

## Examples

To make your socket accessible only to that machine, bind to IP 127.0.0.1 (localhost):

~~~ txt
bind 127.0.0.1
~~~

