+++
title = "root"
description = "*root* simply specifies the root of where to find (zone) files."
weight = 36
tags = ["plugin", "root"]
categories = ["plugin"]
date = "2020-02-06T11:54:58.8775882"
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
