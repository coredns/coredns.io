+++
title = "log"
description = "*log* enables query logging to standard output."
weight = 18
tags = [ "plugin", "log" ]
categories = [ "plugin" ]
date = "2017-12-02T07:46:55.244574"
+++

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

You can further specify the class of responses that get logged:

~~~ txt
log [NAME] [FORMAT] {
    class [success|denial|error|all]
}
~~~

Here `success` `denial` and `error` denotes the class of responses that should be logged. The
classes have the following meaning:

* `success`: successful response
* `denial`: either NXDOMAIN or NODATA (name exists, type does not)
* `error`: SERVFAIL, NOTIMP, REFUSED, etc. Anything that indicates the remote server is not willing to
    resolve the request.
* `all`: the default - nothing is specified.

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
* `{remote}`: client's IP address
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
`{remote} - [{when}] "{type} {class} {name} {proto} {size} {>do} {>bufsize}" {rcode} {>rflags} {rsize} {duration}`
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
