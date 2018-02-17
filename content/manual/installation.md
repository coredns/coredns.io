# Installation

CoreDNS is written in Go, but unless you want to develop plugins are compile CoreDNS your self you
probably don't care. The following sections detail how you can get CoreDNS or install from source.

## Binaries

For every CoreDNS release we provide [pre-compiled
binaries](https://github.com/coredns/coredns/releases/latest) for various operating systems. For
Linux we also provide cross compiled binaries for ARM, PowerPC and other architectures.

## Docker

Also Docker. TODO().

Note that Docker images that are for architectures other than *AMD64*, don't have any certificates
installed. This means if you want to use CoreDNS on ARM and do things like DNS-over-TLS you'll need
to create your own Docker image.

## Source

To compile CoreDNS we assume you have a working Go setup, see various tutorials (/links) if you
don't have that already configured. The Go version that comes with your OS is probably too old to
compile CoreDNS as we require Go 1.9.x at the moment (Feb 2018).

With CoreDNS we try to vendor all our dependencies, but because of [various
reasons](https://github.com/coredns/coredns/issues/1523) (mostly making it
possible for external plugins to compile), we can not vendor *all* our dependencies. Hence to compile
CoreDNS, you still need to `go get`. The `Makefile` we include handles all of these steps. So
compiling CoreDNS boils down to (as of this writing the latest version is 1.0.5):

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

When all done you should end up with a `coredns` executable in the current directory:
~~~
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
