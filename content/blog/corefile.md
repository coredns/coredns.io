+++
date = "2017-07-23T21:28:13Z"
description = "How does the Corefile work?"
tags = ["Corefile", "Documentation"]
title = "Corefile Explained"
author = "miek"
+++

The `Corefile` is CoreDNS's configuration file. It defines:

* What servers listen on what ports.
* For which zone each server is authoritative.
* What middleware is loaded in a server.

To explain more, let take a look at this "Corefile":

~~~ txt
ZONE:[PORT] {
    [MIDDLEWARE]...
}
~~~

* **ZONE** defines the zone this server. The optional **PORT** defaults to 53, if not given (or
  whatever the `-dns.port` flag has as a value.
* **MIDDLEWARE** defines the middleware we want to load. This is optional as well, but as server
  with no middleware will just return SERVFAIL for all queries.

So this is a minimum Corefile:

~~~ txt
. { }
~~~

That defines a server to listen on port 53 and make it authoritative for the root zone *and*
everything below. Let's define another server that is authoritative for `.` (root zone) and load
that:

~~~ txt
. { }
. { }
~~~

This will make CoreDNS exit with an error:
~~~
2017/07/23 20:39:10 cannot serve dns://.:53 - zone already defined for dns://.:53
~~~

Why? Because we already defined a server *on the same port* for this zone. If we change the port
number on the second server and thereby creating *another* server, it is OK:

~~~ txt
.    { }
.:54 { }
~~~

When defining a new zone, you either create a new server, or add it to an existing one - but you can
redefine the middleware for it. Here we define *one* server that handles two zones; that potentially
chain different middleware:

~~~ txt
example.org {
    whoami
}
org {
    whoami
}
~~~

Note that most specific zone wins when a query comes in, so any `example.org` queries are going
through the middleware defined for `example.org` above. The rest is handled by `.`.

## Default Corefile

When CoreDNS starts up and can't find a Corefile, it will load a default one, defined as

~~~ txt
. {
    whoami
}
~~~

## Non Default Protocols

Listening on TLS and for gRPC? Use:

~~~ txt
tls://example.org grpc://example.org {
    # ...
}
~~~

Specifying ports works in the same way, here when listening for gRPC packets.

~~~ txt
grpc://example.org:1443 {
    # ...
}
~~~

## Also See

The Corefile is parsed like a [Caddyfile](https://caddyserver.com/docs/caddyfile).
