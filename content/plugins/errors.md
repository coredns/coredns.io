+++
title = "errors"
description = "*errors* enables error logging."
weight = 10
tags = [ "plugin", "errors" ]
categories = [ "plugin" ]
date = "2017-10-19T06:31:53.688967"
+++

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
