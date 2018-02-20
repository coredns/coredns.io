+++
title = "dump"
description = "*dump* - dumps all incoming queries on standard output."
weight = 10
tags = [  "plugin" , "dump" ]
categories = [ "plugin", "external" ]
date = "2018-01-17T21:57:00+08:00"
repo = "https://github.com/miekg/dump"
home = "https://github.com/miekg/dump/blob/master/README.md"
+++

## Description

*dump* uses the synax from the *log* plugin, and defaults to this format:

~~~
{remote} - [{when}] {>id} {type} {class} {name} {proto} {port}
~~~

So a query will show up as:

~~~
:1 - [17/Jan/2018:20:02:19 +0000] 3644 MX IN example.net. udp 46481
~~~

Note that this is shorter than the default for *log* so you can distinguish between the two outputs.
*log* only logs queries that have seen a response, so this plugin can be used as a debugging aid to
 just dump all incoming queries.

## Syntax

~~~ txt
dump
~~~

## Examples

Dump all queries.

~~~ corefile
. {
    dump
}
~~~
