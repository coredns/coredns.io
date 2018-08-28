+++
title = "log"
description = "*log* enables query logging to standard output."
weight = 19
tags = [ "plugin", "log" ]
categories = [ "plugin" ]
date = "2018-08-28T06:15:01.556838"
+++

## Description

By just using *log* you dump all queries (and parts for the reply) on standard output. Options exist
to tweak the output a little.

Note that for busy servers this will incur a performance hit.

## Syntax

~~~ txt
log
~~~

* With no arguments, a query log entry is written to *stdout* in the common log format for all requests

Or if you want/need slightly more control:

~~~ txt
log [NAME] [FORMAT]
~~~

* `NAME` is the name to match in order to be logged
* `FORMAT` is the log format to use (default is Common Log Format)

You can further specify the classes of responses that get logged:

~~~ txt
log [NAME] [FORMAT] {
    class CLASSES...
}
~~~

* `CLASSES` is a space-separated list of classes of responses that should be logged

The classes of responses have the following meaning:

* `success`: successful response
* `denial`: either NXDOMAIN or NODATA (name exists, type does not)
* `error`: SERVFAIL, NOTIMP, REFUSED, etc. Anything that indicates the remote server is not willing to
    resolve the request.
* `all`: the default - nothing is specified. Using of this class means that all messages will be logged whatever we mix together with "all".

If no class is specified, it defaults to *all*.

## Log Format

You can specify a custom log format with any placeholder values. Log supports both request and
response placeholders.

The following place holders are supported:

* `{type}`: qtype of the request
* `{name}`: qname of the request
* `{class}`: qclass of the request
* `{proto}`: protocol used (tcp or udp)
* `{when}`: time of the query
* `{remote}`: client's IP address, for IPv6 addresses these are enclosed in brackets: `[::1]`
* `{size}`: request size in bytes
* `{port}`: client's port
* `{duration}`: response duration
* `{rcode}`: response RCODE
* `{rsize}`: response size
* `{>rflags}`: response flags, each set flag will be displayed, e.g. "aa, tc". This includes the qr
  bit as well.
* `{>bufsize}`: the EDNS0 buffer size advertised in the query
* `{>do}`: is the EDNS0 DO (DNSSEC OK) bit set in the query
* `{>id}`: query ID
* `{>opcode}`: query OPCODE

The default Common Log Format is:

~~~ txt
`{remote}:{port} - [{when}] {>id} "{type} {class} {name} {proto} {size} {>do} {>bufsize}" {rcode} {>rflags} {rsize} {duration}`
~~~

## Examples

Log all requests to stdout

~~~ corefile
. {
    log
    whoami
}
~~~

Custom log format, for all zones (`.`)

~~~ corefile
. {
    log . "{proto} Request: {name} {type} {>id}"
}
~~~

Only log denials for example.org (and below to a file)

~~~ corefile
. {
    log example.org {
        class denial
    }
}
~~~

Log all queries which were not resolved successfully

~~~ corefile
. {
    log . {
        class denial error
    }
}
~~~

Log all queries on which we did not get errors

~~~ corefile
. {
    log . {
        class denial success
    }
}
~~~

Also the multiple statements can be OR-ed, for example, we can rewrite the above case as following:

~~~ corefile
. {
    log . {
        class denial
        class success
    }
}
~~~
