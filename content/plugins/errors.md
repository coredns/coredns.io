+++
title = "errors"
description = "*errors* enable error logging."
weight = 10
tags = [ "plugin", "errors" ]
categories = [ "plugin" ]
date = "2018-05-24T08:47:52.444723"
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
