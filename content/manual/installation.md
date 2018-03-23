# Installation

CoreDNS is written in Go, but unless you want to develop plugins or compile CoreDNS yourself you
probably don't care. The following sections detail how you can get CoreDNS binaries or install from source.

## Binaries

For every CoreDNS release, we provide [pre-compiled
binaries](https://github.com/coredns/coredns/releases/latest) for various operating systems. For
Linux, we also provide cross compiled binaries for ARM, PowerPC and other architectures.

## Docker

We push every release as Docker images as well. You can find them in the [public Docker
hub](https://hub.docker.com/r/coredns/coredns/) for the CoreDNS organization.

Note that Docker images that are for architectures other than *AMD64* don't have any certificates
installed. This means if you want to use CoreDNS on ARM and do things like DNS-over-TLS, you'll need
to create your own Docker image.

## Source

To compile CoreDNS, we assume you have a working Go setup. See various tutorials if you don't have
that already configured. The Go version that comes with your OS is probably too old to compile
CoreDNS as we require Go 1.9.x at the moment (Feb 2018).

With CoreDNS, we try to vendor all our dependencies, but because of [various
reasons](https://github.com/coredns/coredns/issues/1523) (mostly making it
possible for external plugins to compile), we can not vendor *all* our dependencies. Hence to compile
CoreDNS, you still need to `go get` some packages. The `Makefile` we include handles all of these
steps. So compiling CoreDNS boils down to (as of this writing the latest version is 1.0.5):

~~~ sh
$ export GOPATH=${GOPATH-~/go}
$ mkdir -p $GOPATH/src/github.com/coredns
$ cd $GOPATH/src/github.com/coredns/
$ wget https://github.com/coredns/coredns/archive/v1.0.5.tar.gz
$ tar xvf v1.0.5.tar.gz
$ mv coredns-1.0.5 coredns
$ cd coredns
$ make CHECKS= godeps all
~~~

When all of that is done, you should end up with a `coredns` executable in the current directory:
~~~ sh
$ ./coredns -version
CoreDNS-1.0.5
linux/amd64, go1.9.4,
~~~
The `go1.9.4,` part usually shows a git commit, but as this is a source tar ball, we don't have
this.

## Source from Github

This is mostly the same set of steps:

~~~ sh
$ export GOPATH=${GOPATH-~/go}
$ mkdir -p $GOPATH/src/github.com/coredns
$ cd $GOPATH/src/github.com/coredns/
$ git clone git@github.com:coredns/coredns
$ cd coredns
$ make CHECKS= godeps all
~~~

## Testing

Once you have a `coredns` binary, you can use the `-plugins` flag to list all the compiled plugins.
Without a `Corefile` (See [Configuration](#configuration)) CoreDNS will load the
[*whoami*](/plugins/whoami) that will respond with the IP address and port of the client. So to
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
