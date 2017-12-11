+++
title = "dnstap"
description = "*dnstap* enables logging to dnstap, a flexible, structured binary log format for DNS software: http://dnstap.info."
weight = 8
tags = [ "plugin", "dnstap" ]
categories = [ "plugin" ]
date = "2017-12-11T16:50:50.551637"
+++

There is a buffer, expect at least 13 requests before the server sends its dnstap messages to the socket.

## Syntax

~~~ txt
dnstap SOCKET [full]
~~~

* **SOCKET** is the socket path supplied to the dnstap command line tool.
* `full` to include the wire-format DNS message.

## Examples

Log information about client requests and responses to */tmp/dnstap.sock*.

~~~ txt
dnstap /tmp/dnstap.sock
~~~

Log information including the wire-format DNS message about client requests and responses to */tmp/dnstap.sock*.

~~~ txt
dnstap unix:///tmp/dnstap.sock full
~~~

Log to a remote endpoint.

~~~ txt
dnstap tcp://127.0.0.1:6000 full
~~~

## Dnstap command line tool

~~~ sh
% go get github.com/dnstap/golang-dnstap
% cd $GOPATH/src/github.com/dnstap/golang-dnstap/dnstap
% go build
% ./dnstap
~~~

The following command listens on the given socket and decodes messages to stdout.

~~~ sh
% dnstap -u /tmp/dnstap.sock
~~~

The following command listens on the given socket and saves message payloads to a binary dnstap-format log file.

~~~ sh
% dnstap -u /tmp/dnstap.sock -w /tmp/test.dnstap
~~~

Listen for dnstap messages on port 6000.

~~~ sh
% dnstap -l 127.0.0.1:6000
~~~
