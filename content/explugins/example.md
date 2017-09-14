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

# example

The example plugin prints "example" on every query received. It can be used as documentation for
writing external plugin and to test if external plugin compiles with CoreDNS.

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

## How to Enable

Follow [these](/2017/07/25/compile-time-enabling-or-disabling-plugin/) steps,
*example* should be put relatively early in the plugin chain.
