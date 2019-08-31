+++
title = "erratic"
description = "*erratic* a plugin useful for testing client behavior."
weight = 13
tags = [ "plugin", "erratic" ]
categories = [ "plugin" ]
date = "2019-08-31T08:36:24.153297"
+++

## Description

*erratic* returns a static response to all queries, but the responses can be delayed, dropped or truncated.
The *erratic* plugin will respond to every A or AAAA query. For any other type it will return
a SERVFAIL response. The reply for A will return 192.0.2.53 (see [RFC
5737](https://tools.ietf.org/html/rfc5737),
for AAAA it returns 2001:DB8::53 (see [RFC 3849](https://tools.ietf.org/html/rfc3849)) and for an
AXFR request it will respond with a small zone transfer.

*erratic* can also be used in conjunction with the *autopath* plugin. This is mostly to aid in
testing.

## Syntax

~~~ txt
erratic {
    drop [AMOUNT]
    truncate [AMOUNT]
    delay [AMOUNT [DURATION]]
}
~~~

* `drop`: drop 1 per **AMOUNT** of queries, the default is 2.
* `truncate`: truncate 1 per **AMOUNT** of queries, the default is 2.
* `delay`: delay 1 per **AMOUNT** of queries for **DURATION**, the default for **AMOUNT** is 2 and
  the default for **DURATION** is 100ms.

In case of a zone transfer and truncate the final SOA record *isn't* added to the response.

## Ready

This plugin reports readiness to the ready plugin.

## Examples

~~~ corefile
. {
    erratic {
        drop 3
    }
}
~~~

Or even shorter if the defaults suits you. Note this only drops queries, it does not delay them.

~~~ corefile
. {
    erratic
}
~~~

Delay 1 in 3 queries for 50ms

~~~ corefile
. {
    erratic {
        delay 3 50ms
    }
}
~~~

Delay 1 in 3 and truncate 1 in 5.

~~~ corefile
. {
    erratic {
        delay 3 5ms
        truncate 5
    }
}
~~~

Drop every second query.

~~~ corefile
. {
    erratic {
        drop 2
        truncate 2
    }
}
~~~

## Also See

[RFC 3849](https://tools.ietf.org/html/rfc3849) and
[RFC 5737](https://tools.ietf.org/html/rfc5737).
