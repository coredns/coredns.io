+++
title = "any"
description = "*any* - give a minimal response to ANY queries."
weight = 10
tags = [  "plugin" , "any" ]
categories = [ "plugin", "external" ]
date = "2019-03-25T20:17:00+08:00"
repo = "https://github.com/coredns/any"
home = "https://github.com/coredns/any/blob/master/README.md"
+++

## Description

*any* basically blocks ANY queries by responding to it with a short HINFO reply.
See [RFC 8482](https://tools.ietf.org/html/rfc8482) for details.

## Syntax

~~~ txt
any
~~~

## Examples

~~~ corefile
. {
    whoami
    any
}
~~~

## Also See

[RFC 8482](https://tools.ietf.org/html/rfc8482).
