+++
title = "errors"
description = "*errors* enables error logging."
weight = 10
tags = [ "middleware", "errors" ]
categories = [ "middleware" ]
date = "2017-09-10T18:11:52.763391"
+++

Any errors encountered during the query processing will be printed on standard output.

## Syntax

~~~
errors [FILE]
~~~

* **FILE** is the log file to create (or append to). The *only* valid name for **FILE** is *stdout*

## Examples

Log errors to *stdout*.

~~~
errors
~~~

Log errors to *stdout*.

~~~
errors stdout
~~~
