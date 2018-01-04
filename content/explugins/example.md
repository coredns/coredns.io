+++
title = "example"
description = "*example* prints 'example' on every query received."
weight = 10
tags = [  "plugin" , "example" ]
categories = [ "plugin", "external" ]
date = "2017-07-25T21:57:00+08:00"
repo = "https://github.com/coredns/example/"
home = "https://github.com/coredns/example/blob/master/README.md"
+++

## Description

The example middleware prints "example" on every query received. It can be used as documentation for
writing external middleware and to test if external middleware compiles with CoreDNS.

## Syntax

~~~ txt
example
~~~

## Examples

```
example.com {
  file example.com.db {
    upstream 8.8.8.8
  }
  example
}
```
