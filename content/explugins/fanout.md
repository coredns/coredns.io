+++
title = "fanout"
description = "*fanout* - parallel proxying DNS messages to upstream resolvers."
weight = 10
tags = [  "plugin" , "fanout" ]
categories = [ "plugin", "external" ]
date = "2024-09-03T22:00:00+08:00"
repo = "https://github.com/networkservicemesh/fanout"
home = "https://github.com/networkservicemesh/fanout/README.md"
+++

## Description

Each incoming DNS query that hits the CoreDNS fanout plugin will be replicated in parallel to each listed IP (i.e. the DNS servers). The first non-negative response from any of the queried DNS Servers will be forwarded as a response to the application's DNS request.

## Syntax

* `tls` **CERT** **KEY** **CA** define the TLS properties for TLS connection. From 0 to 3 arguments can be
  provided with the meaning as described below
  * `tls` - no client authentication is used, and the system CAs are used to verify the server certificate
  * `tls` **CA** - no client authentication is used, and the file CA is used to verify the server certificate
  * `tls` **CERT** **KEY** - client authentication is used with the specified cert/key pair.
    The server certificate is verified with the system CAs
  * `tls` **CERT** **KEY**  **CA** - client authentication is used with the specified cert/key pair.
    The server certificate is verified using the specified CA file
* `tls_servername` **NAME** allows you to set a server name in the TLS configuration; for instance 9.9.9.9
  needs this to be set to `dns.quad9.net`. Multiple upstreams are still allowed in this scenario,
  but they have to use the same `tls_servername`. E.g. mixing 9.9.9.9 (QuadDNS) with 1.1.1.1
  (Cloudflare) will not work.

* `worker-count` is the number of parallel queries per request. By default equals to count of IP list. Use this only for reducing parallel queries per request.
* `policy` - specifies the policy of DNS server selection mechanism. The default is `sequential`.
  * `sequential` - select DNS servers one-by-one based on its order
  * `weighted-random` - select DNS servers randomly based on `weighted-random-server-count` and `weighted-random-load-factor` params.
* `weighted-random-server-count` is the number of DNS servers to be requested. Equals to the number of specified IPs by default. Used only with the `weighted-random` policy.
* `weighted-random-load-factor` - the probability of selecting a server. This is specified in the order of the list of IP addresses and takes values between 1 and 100. By default, all servers have an equal probability of 100. Used only with the `weighted-random` policy.
* `network` is a specific network protocol. Could be `tcp`, `udp`, `tcp-tls`.
* `except` is a list is a space-separated list of domains to exclude from proxying.
* `except-file` is the path to file with line-separated list of domains to exclude from proxying.
* `attempt-count` is the number of attempts to connect to upstream servers that are needed before considering an upstream to be down. If 0, the upstream will never be marked as down and request will be finished by `timeout`. Default is `3`.
* `timeout` is the timeout of request. After this period, attempts to receive a response from the upstream servers will be stopped. Default is `30s`.
* `race` gives priority to the first result, whether it is negative or not, as long as it is a standard DNS result.
## Metrics

If monitoring is enabled (via the *prometheus* plugin) then the following metric are exported:

* `coredns_fanout_request_duration_seconds{to}` - duration per upstream interaction.
* `coredns_fanout_request_count_total{to}` - query count per upstream.
* `coredns_fanout_response_rcode_count_total{to, rcode}` - count of RCODEs per upstream.

Where `to` is one of the upstream servers (**TO** from the config), `rcode` is the returned RCODE
from the upstream.

## Examples
Proxy all requests within `example.org.` to a nameservers running on a different ports.  The first positive response from a proxy will be provided as the result.

~~~ corefile
example.org {
    fanout . 127.0.0.1:9005 127.0.0.1:9006 127.0.0.1:9007 127.0.0.1:9008
}
~~~

Sends parallel requests between three resolvers, one of which has a IPv6 address via TCP. The first response from proxy will be provided as the result.

~~~ corefile
. {
    fanout . 10.0.0.10:53 10.0.0.11:1053 [2003::1]:53 {
        network TCP
    }
}
~~~

Proxying everything except requests to `example.org`

~~~ corefile
. {
    fanout . 10.0.0.10:1234 {
        except example.org
    }
}
~~~

Proxy everything except `example.org` using the host's `resolv.conf`'s nameservers:

~~~ corefile
. {
    fanout . /etc/resolv.conf {
        except example.org
    }
}
~~~

Proxy all requests to 9.9.9.9 using the DNS-over-TLS protocol.
Note the `tls-server` is mandatory if you want a working setup, as 9.9.9.9 can't be
used in the TLS negotiation.

~~~ corefile
. {
    fanout . tls://9.9.9.9 {
       tls-server dns.quad9.net
    }
}
~~~

Sends parallel requests between five resolvers via UDP uses two workers and without attempting to reconnect. The first positive response from a proxy will be provided as the result.
~~~ corefile
. {
    fanout . 10.0.0.10:53 10.0.0.11:53 10.0.0.12:53 10.0.0.13:1053 10.0.0.14:1053 {
        worker-count 2
    }
}
~~~

Multiple upstream servers are configured but one of them is down, query a `non-existent` domain.
If `race` is enable, we will get `NXDOMAIN` result quickly, otherwise we will get `"connection timed out"` result in a few seconds.
~~~ corefile
. {
    fanout . 10.0.0.10:53 10.0.0.11:53 10.0.0.12:53 10.0.0.13:1053 10.0.0.14:1053 {
        race
    }
}
~~~

Sends parallel requests between two randomly selected resolvers. Note, that `127.0.0.1:9007` would be selected more frequently as it has the highest `weighted-random-load-factor`.
~~~ corefile
example.org {
    fanout . 127.0.0.1:9005 127.0.0.1:9006 127.0.0.1:9007 {
      policy weighted-random
      weighted-random-server-count 2
      weighted-random-load-factor 50 70 100
    }
}
~~~

Sends parallel requests between three resolver sequentially (default mode).
~~~ corefile
example.org {
    fanout . 127.0.0.1:9005 127.0.0.1:9006 127.0.0.1:9007 {
        policy sequential
    }
}
~~~

## Also See

See the [fanout](https://github.com/networkservicemesh/fanout).