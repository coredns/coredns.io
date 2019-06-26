+++
title = "Logging with dnstap"
description = "*dnstap* is a flexible, structured binary log format for DNS software."
tags = ["dnstap", "log", "plugin"]
date = "2017-08-03T10:25:28+02:00"
author = "varyoo"
+++

[^dnstap]: http://dnstap.info

**dnstap** is a flexible, structured binary log format for DNS software[^dnstap].
It uses [Protocol Buffers](https://developers.google.com/protocol-buffers/) to encode events that
occur inside DNS software in an implementation-neutral format.

dnstap can encode any DNS message exchanged by the server, along with information about the remote
computer (IP address, port) and time. It includes client queries and responses, but also proxied
requests or other information requested from other name servers.

This example shows output from the dnstap command-line tool to get an idea of the kind of
information that dnstap can provide:

~~~ text
type: MESSAGE
message:
  type: CLIENT_RESPONSE
  socket_family: INET
  socket_protocol: UDP
  query_address: 127.0.0.1
  query_port: 47969
  response_message: |
    ;; opcode: QUERY, status: NOERROR, id: 47163
    ;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

    ;; QUESTION SECTION:
    ;example.org.       IN       A

    ;; ANSWER SECTION:
    example.org.        86339   IN      A       93.184.216.34
~~~

A [*dnstap* plugin]({{< relref "/plugins/dnstap.md" >}}) has been added in [CoreDNS-010]({{< relref "/blog/coredns-010.md" >}}).
Currently it can only log client level messages. Logging for additional types of exchanges is being implemented.

The *dnstap* plugin is used in combination with the **dnstap** command-line tool.
They use a socket to communicate: the plugin will send the logs as long as the tool is listening.

To start with the *dnstap* plugin add it to the Corefile in a server block:

~~~ text
dnstap /tmp/dnstap.sock full
~~~

With the `full` option given to the *dnstap* plugin you will also include the full (binary) data
of the DNS message.
Now you can use the *dnstap* tool to read from the socket where CoreDNS writes to.

~~~ text
$ dnstap -u /tmp/dnstap.sock
~~~

Or listen on the *dnstap* socket and store message payloads to a binary *dnstap*-format log file:

~~~ text
$ dnstap -u /tmp/dnstap.sock -w /tmp/july.dnstap
~~~

And then read back July's logs in the YAML-format:

~~~ text
$ dnstap -r /tmp/july.dnstap -y
~~~
