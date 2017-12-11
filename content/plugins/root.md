+++
title = "root"
description = "*root* simply specifies the root of where to find (zone) files."
weight = 25
tags = [ "plugin", "root" ]
categories = [ "plugin" ]
date = "2017-12-11T16:50:50.559672"
+++

The default root is the current working directory of CoreDNS. A relative root path is relative to
the current working directory.

## Syntax

~~~ txt
root PATH
~~~

**PATH** is the directory to set as CoreDNS' root.

## Examples

Serve zone data (when the *file* plugin is used) from `/etc/coredns/zones`:

~~~ corefile
. {
    root /etc/coredns/zones
}
~~~
