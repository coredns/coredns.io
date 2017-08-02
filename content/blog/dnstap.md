+++
title = "Logging with dnstap"
description = "dnstap is a flexible, structured binary log format for DNS software."
tags = ["dnstap", "log", "middleware"]
draft = true
date = "2017-08-01T16:25:28+02:00"
author = "varyoo"
+++

*dnstap* is a flexible, structured binary log format for DNS software.
It uses [Protocol Buffers](https://developers.google.com/protocol-buffers/) to encode events that occur inside DNS software in an implementation-neutral format.

*dnstap* can encode any DNS message exchanged by the server, along with information about the remote computer (IP address, port) and time.
It includes client queries and responses, but also proxied requests or information requested from other name servers.

Check out this example output from the *dnstap* command-line tool to get an idea of the kind of information that *dnstap* can encode:

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

A [*dnstap* middleware] has been added in [CoreDNS-010]({{< relref "blog/coredns-010.md" >}}).
Currently it can only log client level messages. Logging additional type of exchanges is being experimented.

The [*dnstap* middleware] is used in combination with the *dnstap* command-line tool.
They use a socket to communicate:
the middleware will send the logs as long as the CLI tool is listening.

# Examples

Add *dnstap* to the *Corefile*:

~~~ text
dnstap /tmp/dnstap.sock full
~~~

Listen on the *dnstap* socket and write message payloads to *stdout*:

~~~ text
$ dnstap -u /tmp/dnstap.sock
~~~

Listen on the *dnstap* socket and store message payloads to a binary *dnstap*-format log file:

~~~ text
$ dnstap -u /tmp/dnstap.sock -w /tmp/july.dnstap
~~~

Read July's logs in the YAML-format:

~~~ text
$ dnstap -r /tmp/july.dnstap -y
~~~

[*dnstap* middleware]: {{< relref "middleware/dnstap.md" >}}
