+++
title = "any"
description = "*any* give a minimal response to ANY queries."
weight = 1
tags = [ "plugin", "any" ]
categories = [ "plugin" ]
date = "2019-07-03T18:33:28.050440"
+++

## Description

*any* basically blocks ANY queries by responding to them with a short HINFO reply. See [RFC
8482](https://tools.ietf.org/html/rfc8482) for details.

## Syntax

~~~ txt
any
~~~

## Examples

~~~ corefile
example.org {
    whoami
    any
}
~~~

A `dig +nocmd ANY example.org +noall +answer` now returns:

~~~ txt
example.org.  8482	IN	HINFO	"ANY obsoleted" "See RFC 8482"
~~~

## Also See

[RFC 8482](https://tools.ietf.org/html/rfc8482).
