+++
title = "root"
description = "*root* simply specifies the root of where to find (zone) files."
weight = 33
tags = [ "plugin", "root" ]
categories = [ "plugin" ]
date = "2019-06-26T12:57:30.985229"
+++

## Description

The default root is the current working directory of CoreDNS. The *root* plugin allows you to change
this. A relative root path is relative to the current working directory.

This plugin can only be used once per Server Block.

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
