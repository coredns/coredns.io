+++
title = "debug"
description = "*debug* disables the automatic recovery upon a crash so that you'll get a nice stack trace."
weight = 6
tags = [ "plugin", "debug" ]
categories = [ "plugin" ]
date = "2018-02-22T08:55:16.399368"
+++

## Description

Normally CoreDNS will recover from panics, using *debug* inhibits this. The main use of *debug* is
to help testing.

Note that the *errors* plugin (if loaded) will also set a `recover` negating this setting. 

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
