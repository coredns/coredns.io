+++
date = "2017-07-23T21:28:13Z"
description = "How does the Corefile work?"
tags = ["Corefile", "Documentation"]
title = "Corefile Explained"
author = "miek"
+++

The `Corefile` is CoreDNS's configuration file. It defines:

* What servers listen on what ports and which protocol.
* For which zone each server is authoritative.
* Which plugins are loaded in a server.

To explain more, let take a look at this "Corefile":

~~~ txt
ZONE:[PORT] {
    [PLUGIN]...
}
~~~

* **ZONE** defines the zone this server. The optional **PORT** defaults to 53, or
  the value of the `-dns.port` flag.
* **PLUGIN** defines the [plugin(s)](/plugins) we want to load. This is optional as well, but as server
  with no plugins will just return SERVFAIL for all queries. Each plugin can have a number of
  *properties* than can have *arguments*

I.e., in the next example:

The **ZONE** is root zone `.`, the **PLUGIN** is `chaos`. The [*chaos* plugin](/plugins/chaos) does
not have any properties, but it does take an *argument*: `CoreDNS-001`. This text is returned on
a CH class query: `dig CH txt version.bind @localhost`

~~~ corefile
. {
   chaos CoreDNS-001
}
~~~

If CoreDNS can't find a `Corefile` to load is loads the following builtin one that loads the
[*whoami* plugin](/plugin/whoami):

~~~ corefile
. {
    whoami
}
~~~

## Servers

This is the most minimal Corefile:

~~~ corefile
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

When defining a new zone, you either create a new server, or add it to an existing one. Here we
define *one* server that handles two zones; that potentially chain different plugin:

~~~ corefile
example.org {
    whoami
}
org {
    whoami
}
~~~

Note that most specific zone wins when a query comes in, so any `example.org` queries are going
through the server defined for `example.org` above. The queries for `.org` are going to the other 
server.

## Reverse Zones

Normally when you want to serve a reverse zone you'll have to say something:

~~~ corefile
0.0.10.in-addr.arpa {
    whoami
}
~~~

To make this easier CoreDNS just allows you to say:

~~~ corefile
10.0.0.0/24 {
    whoami
}
~~~

This also works for CIDR (in the 1.0.0 release) zones:

~~~ corefile
10.0.0.0/27 {
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

The Corefile is parsed like a [Caddyfile](https://caddyserver.com/docs/caddyfile). We support
everything that is described on that page, for instance the use of [environment
variables](https://caddyserver.com/docs/caddyfile#env).

Other interesting plugins that are helpful in Corefiles are:
[*import*](/plugins/import)
[*startup*](/plugins/startup) and
[*shutdown*](/plugins/shutdown).
