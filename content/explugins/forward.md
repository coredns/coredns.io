+++
title = "forward"
description = "*forward* facilitates proxying DNS messages to upstream resolvers."
weight = 10
tags = [  "plugin" , "forward" ]
categories = [ "plugin", "external" ]
date = "2017-10-10T18:25:19+01:00"
repo = "https://github.com/coredns/forward"
home = "https://github.com/coredns/forward"
+++

## Description

*forward* facilitates proxying DNS messages to upstream resolvers.

The *forward* plugin is generally faster (~30%) than *proxy* as it re-uses already openened sockets
to the upstreams. It supports UDP, TCP and DNS-over-TLS and uses inband healthchecking that is
enabled by default.

## Syntax

In its most basic form, a simple forwarder uses this syntax:

~~~
forward FROM TO...
~~~

* **FROM** is the base domain to match for the request to be forwarded.
* **TO...** are the destination endpoints to forward to. The **TO** syntax allows you to specify
  a protocol, `tls://9.9.9.9` or `dns://` for plain DNS.

The health checks are done every *0.5s*. After *two* failed checks the upstream is considered
unhealthy. The health checks use a recursive DNS query (`. IN NS`) to get upstream health. Any
response that is not an error (REFUSED, NOTIMPL, SERVFAIL, etc) is taken as a healthy upstream. The
health check uses the same protocol as specific in the **TO**. On startup each upstream is marked
unhealthy until it passes a healthcheck.

Multiple upstreams are randomized on first use. When a healthy proxy returns an error during the
exchange the next upstream in the list is tried.

Extra knobs are available with an expanded syntax:

~~~
forward FROM TO... {
    except IGNORED_NAMES...
    force_tcp
    health_check DURATION
    expire DURATION
    max_fails INTEGER
    tls CERT KEY CA
    tls_servername NAME
}
~~~

* **FROM** and **TO...** as above.
* **IGNORED_NAMES** in `except` is a space-separated list of domains to exclude from forwarding.
  Requests that match none of these names will be passed through.
* `force_tcp`, use TCP even when the request comes in over UDP.
* `health_checks`, use a different **DURATION** for health checking, the default duration is 500ms.
* `max_fails` is the number of subsequent failed health checks that are needed before considering
  a backend to be down. If 0, the backend will never be marked as down. Default is 2.
* `expire` **DURATION**, expire connections after this time, the default is 10s.
* `tls` **CERT** **KEY** **CA** define the TLS properties for TLS; if you leave this out the
  system's configuration will be used.
* `tls_servername` **NAME** allows you to set a server name in the TLS configuration; for instance 9.9.9.9
  needs this to be set to `dns.quad9.net`.

The upstream selection is done via round robin. If the socket for this client isn't known *forward*
will randomly choose one. If this turns out to be unhealthy, the next one is tried.

Also note the TLS config is "global" for the whole forwarding proxy if you need a different
`tls-name` for different upstreams you're out of luck.

## Metrics

If monitoring is enabled (via the *prometheus* directive) then the following metric are exported:

* `coredns_forward_request_duration_seconds{to}` - duration per upstream interaction.
* `coredns_forward_request_count_total{to}` - query count per upstream.
* `coredns_forward_response_rcode_total{to, rcode}` - count of RCODEs per upstream.
* `coredns_forward_healthcheck_failure_count_total{to}` - number of failed healthchecks per upstream.
* `coredns_forward_socket_count_total{to}` - number of cached sockets per upstream.

Where `to` is one of the upstream servers (**TO** from the config), `proto` is the protocol used by
the incoming query ("tcp" or "udp"), and family the transport family ("1" for IPv4, and "2" for
IPv6).

## Examples

Proxy all requests within example.org. to a nameserver running on a different port:

~~~ corefile
example.org {
    forward . 127.0.0.1:9005
}
~~~

Load-balance all requests between three resolvers:

~~~ corefile
. {
    forward . 10.0.0.10:53 10.0.0.11:1053 10.0.0.12
}
~~~

Forward everything except requests to `example.org`

~~~ corefile
. {
    forward . 10.0.0.10:1234 {
        except example.org
    }
}
~~~

Proxy everything except `example.org` using the host's `resolv.conf`'s nameservers:

~~~ corefile
. {
    forward . /etc/resolv.conf {
        except example.org
    }
}
~~~

Forward to a IPv6 host:

~~~ corefile
. {
    forward . [::1]:1053
}
~~~

Proxy all requests to 9.9.9.9 using the DNS-over-TLS protocol, and cache every answer for up to 30
seconds.

~~~ corefile
. {
    forward . tls://9.9.9.9 {
       tls_servername dns.quad9.net
       health_check 5s
    }
    cache 30
}
~~~
