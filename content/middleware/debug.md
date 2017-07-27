+++
title = "debug"
description = "*debug* disables the automatic recovery upon a CoreDNS crash so that you'll get a nice stack trace."
weight = 5
tags = [ "middleware", "debug" ]
categories = [ "middleware" ]
date = "2017-07-27T12:53:47.833872"
+++

Note that the *errors* middleware (if loaded) will also set a `recover` negating this setting.
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
