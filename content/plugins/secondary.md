+++
title = "secondary"
description = "*secondary* enables serving a zone retrieved from a primary server."
weight = 37
tags = [ "plugin", "secondary" ]
categories = [ "plugin" ]
date = "2019-09-27T10:37:57.665504"
+++

## Description

With *secondary* you can transfer (via AXFR) a zone from another server. The retrieved zone is
*not committed* to disk (a violation of the RFC). This means restarting CoreDNS will cause it to
 retrieve all secondary zones.

~~~
secondary [ZONES...]
~~~

* **ZONES** zones it should be authoritative for. If empty, the zones from the configuration block
    are used. Note that without a remote address to *get* the zone from, the above is not that useful.

A working syntax would be:

~~~
secondary [zones...] {
    transfer from ADDRESS
    transfer to ADDRESS
}
~~~

* `transfer from` specifies from which address to fetch the zone. It can be specified multiple times;
    if one does not work, another will be tried.
* `transfer to` can be enabled to allow this secondary zone to be transferred again.

When a zone is due to be refreshed (Refresh timer fires) a random jitter of 5 seconds is
applied, before fetching. In the case of retry this will be 2 seconds. If there are any errors
during the transfer the transfer fails; this will be logged.

## Examples

Transfer `example.org` from 10.0.1.1, and if that fails try 10.1.2.1.

~~~ corefile
example.org {
    secondary {
        transfer from 10.0.1.1
        transfer from 10.1.2.1
    }
}
~~~

Or re-export the retrieved zone to other secondaries.

~~~ corefile
. {
    secondary example.net {
        transfer from 10.1.2.1
        transfer to *
    }
}
~~~

## Bugs

Only AXFR is supported and the retrieved zone is not committed to disk.
