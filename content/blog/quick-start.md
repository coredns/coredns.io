+++
date = "2017-07-24T07:37:13Z"
description = "Quick Start Guide"
tags = ["Quick", "Start", "Documentation"]
title = "Quick Start"
author = "miek"
+++

First get CoreDNS, either

* *Download the latest* release from [github](https://github.com/coredns/coredns/releases), unpack
  it. You should now have a "coredns" executable.

* *Compile from git* by getting the source code from [github](https://github.com/coredns/coredns).
  Change directory to `coredns` and:

  * `go get` - to get a few dependencies, the other ones are vendored
  * `go build`

  You should now have a "coredns" executable.

* *Get the Docker container* from [docker hub](https://hub.docker.com/r/coredns/coredns/).

If you want to use CoreDNS in Kubernetes, please check [this post about SD with the *kuberneters*
middleware](/2017/03/01/coredns-for-kubernetes-service-discovery-take-2/).

The remainder of this quick start will focus and two different use cases

1. Using CoreDNS to serve zone files. Optionally signing the zones as well.
2. Using CoreDNS as a forwarding proxy.

## Serving from files

When serving from zone files you will want to use the *file* middleware. Let's start with the zone
`example.org.` and zonefile we want to serve from:

Create a file names `/etc/coredns/zones/example.org` with the following content:

~~~ dns
$ORIGIN example.org.
@	3600 IN	SOA sns.dns.icann.org. noc.dns.icann.org. (
				2017042745 ; serial
				7200       ; refresh (2 hours)
				3600       ; retry (1 hour)
				1209600    ; expire (2 weeks)
				3600       ; minimum (1 hour)
				)

    3600 IN NS a.iana-servers.net.
	3600 IN NS b.iana-servers.net.

www     IN A     127.0.0.1
        IN AAAA  ::1
~~~

Create a Corefile, `/etc/coredns/Corefile`, with:

~~~ txt
example.org {
    file /etc/coredns/zones/example.org
    prometheus     # enable metrics
    errors stdout  # show errors
    log stdout     # show query logs
}
~~~

Note: you can put the last 3 directives in another file say "/etc/coredns/zone-default" and use
an `import`: `import /etc/coredns/zone-default`.

Start CoreDNS on a non-standard port to check if everything is correct: `coredns -conf Corefile
-dns.port 1053` and send it a query with [dig](https://en.wikipedia.org/wiki/Dig_(command)):
~~~
% dig -p 1053 @localhost AAAA www.example.org +noall +answer

www.example.org.	3600	IN	AAAA	::1
~~~

As we've enabled logging the query should be show up there as well:
~~~ txt
::1 - [24/Jul/2017:10:10:44 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 121 133.449µs
~~~

From here you can enable CoreDNS to run on port 53 and have it start from systemd (when on Linux),
see [the deployment repo](https://github.com/coredns/deployment) for example scripts.
Or read more about the [*file* middleware](/middleware/file/).

## CoreDNS as proxy

Another way of running CoreDNS is as a proxy, for instance sending DNS request to Google using
HTTPS. Create a Corefile with:

~~~ txt
. {
    proxy . 8.8.8.8:53 {
        protocol https_google
    }
    prometheus
    errors stdout
    log stdout
}
~~~

Start CoreDNS, just like above and send it a few queries. CoreDNS should logs those, in this case:
~~~
::1 - [24/Jul/2017:10:44:15 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 76 83.396955ms
::1 - [24/Jul/2017:10:44:17 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 76 14.030914ms
::1 - [24/Jul/2017:10:44:19 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 76 13.286384ms
~~~

If you look at the time each query took, we measure those in "ms", so it's quite slow. Let's add
caching, but enable the *caching* middleware. Just add the word "cache" to the Corefile and reload
CoreDNS: `kill -SIGUSR1 <pid_of_coredns>`. And query again:

~~~
::1 - [24/Jul/2017:11:33:54 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 76 43.469743ms
::1 - [24/Jul/2017:11:33:55 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 73 133.073µs
~~~

133 µs. That sounds better.

Read more about the [*cache*](/middleware/cache) and [*proxy*](/middleware/proxy) middleware on
this website. And find all other documentation [here](/tags/documentation).
