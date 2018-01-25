+++
title = "errors"
description = "*errors* enable error logging."
weight = 10
tags = [ "plugin", "errors" ]
categories = [ "plugin" ]
date = "2018-01-25T23:18:26.001977"
+++

## Description

Any errors encountered during the query processing will be printed to standard output.

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
