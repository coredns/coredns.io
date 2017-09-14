+++
title = "debug"
description = "*debug* disables the automatic recovery upon a CoreDNS crash so that you'll get a nice stack trace."
weight = 6
tags = [ "plugin", "debug" ]
categories = [ "plugin" ]
date = "2017-09-10T18:11:52.762769"
+++

Note that the *errors* plugin (if loaded) will also set a `recover` negating this setting.
The main use of *debug* is to help testing.

## Syntax

~~~ txt
debug
~~~

## Examples

Disable CoreDNS' ability to recover from crashes:

~~~ txt
debug
~~~
