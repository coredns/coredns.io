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

*dnstap* can log any DNS message exchanged by the server, along with information about the remote computer (IP address, port) and time.
It includes client queries and responses, but also proxied requests or information requested from other name servers.

A [*dnstap* middleware] has been added in [CoreDNS-010]({{< relref "blog/coredns-010.md" >}}).
Currently it can only log client level messages. Logging additional type of exchanges is being experimented.

Check out this example output from the *dnstap* command-line tool to get an idea of the kind of information that *dnstap* logs:

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

The [*dnstap* middleware] is used in combinaison with the *dnstap* command-line tool.
The latest listen the same socket the other is logging to.

# Quick start

1. Listen on the *dnstap* socket:

    ~~~ text
    $ dnstap -u /tmp/dnstap.sock
    ~~~

2. Add *dnstap* to the *Corefile*:

    ~~~ text
    dnstap /tmp/dnstap.sock full
    ~~~

[*dnstap* middleware]: {{< relref "middleware/dnstap.md" >}}
