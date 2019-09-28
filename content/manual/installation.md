# Installation

CoreDNS is written in Go, but unless you want to develop plugins or compile CoreDNS yourself, you
probably don't care. The following sections detail how you can get CoreDNS binaries or install from
source.

## Binaries

For every CoreDNS release, we provide [pre-compiled
binaries](https://github.com/coredns/coredns/releases/latest) for various operating systems. For
Linux, we also provide cross-compiled binaries for ARM, PowerPC and other architectures.

## Docker

We push every release as Docker images as well. You can find them in the [public Docker
hub](https://hub.docker.com/r/coredns/coredns/) for the CoreDNS organization. This Docker image is
basically *scratch* + CoreDNS + TLS certificates (for DoT, DoH, and gRPC).

## Source

To compile CoreDNS, we assume you have a working Go setup. See various tutorials if you don't have
that already configured. CoreDNS is using Go modules for its dependency management.

The most current document on how to compile things is kept in the [coredns
source](https://github.com/coredns/coredns#compilation-from-source).

## Testing

Once you have a `coredns` binary, you can use the `-plugins` flag to list all the compiled plugins.
Without a `Corefile` (See [Configuration](#configuration)) CoreDNS will load the
[*whoami*](/plugins/whoami) plugin that will respond with the IP address and port of the client. So to
test, we start CoreDNS to run on port 1053 and send it a query using `dig`:

~~~ sh
$ ./coredns -dns.port=1053
.:1053
2018/02/20 10:40:44 [INFO] CoreDNS-1.0.5
2018/02/20 10:40:44 [INFO] linux/amd64, go1.10,
CoreDNS-1.0.5
linux/amd64, go1.10,
~~~

And from a different terminal window, a `dig` should return something similar to this:

~~~ sh
$ dig @localhost -p 1053 a whoami.example.org

;; QUESTION SECTION:
;whoami.example.org.		IN	A

;; ADDITIONAL SECTION:
whoami.example.org.	0	IN	AAAA	::1
_udp.whoami.example.org. 0	IN	SRV	0 0 39368 .
~~~

The [next section](#configuration) will show how to enable more interesting plugins.
