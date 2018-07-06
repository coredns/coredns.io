+++
title = "errors"
description = "*errors* enable error logging."
weight = 10
tags = [ "plugin", "errors" ]
categories = [ "plugin" ]
date = "2018-07-06T10:27:55.911392"
+++

## Description

Any errors encountered during the query processing will be printed to standard output.

This plugin can only be used once per Server Block.

## Syntax

~~~
errors
~~~

## Examples

Use the *whoami* to respond to queries and Log errors to standard output.

~~~ corefile
. {
    whoami
    errors
}
~~~
