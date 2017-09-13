+++
title = "root"
description = "*root* simply specifies the root of where CoreDNS finds (e.g.) zone files."
weight = 24
tags = [ "middleware", "root" ]
categories = [ "middleware" ]
date = "2017-09-10T18:11:52.766138"
+++

The default root is the current working directory of CoreDNS. A relative root path is relative to
the current working directory.

## Syntax

~~~ txt
root PATH
~~~

**PATH** is the directory to set as CoreDNS' root.

## Examples

Serve zone data (when the *file* middleware is used) from `/etc/coredns/zones`:

~~~ txt
root /etc/coredns/zones
~~~
