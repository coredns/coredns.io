+++
title = "debug"
description = "*debug* disables the automatic recovery upon a crash so that you'll get a nice stack trace."
weight = 6
tags = [ "plugin", "debug" ]
categories = [ "plugin" ]
date = "2017-12-11T16:50:50.551246"
+++

Note that the *errors* plugin (if loaded) will also set a `recover` negating this setting. The main
use of *debug* is to help testing.

## Syntax

~~~ txt
debug
~~~

## Examples

Disable the ability to recover from crashes:

~~~ corefile
. {
    debug
}
~~~
