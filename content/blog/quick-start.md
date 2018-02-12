+++
date = "2017-07-24T07:37:13Z"
description = "Quick Start Guide."
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

* *Get the Docker image* from [docker hub](https://hub.docker.com/r/coredns/coredns/).

If you want to use CoreDNS in Kubernetes, please check [this post about SD with the *kuberneters*
plugin](/2017/03/01/coredns-for-kubernetes-service-discovery-take-2/).

The remainder of this quick start will focus and two different use cases

1. Using CoreDNS to serve zone files. Optionally signing the zones as well.
2. Using CoreDNS as a forwarding proxy.

CoreDNS is configured via a configuration file that it typically called
[Corefile](https://coredns.io/2017/07/23/corefile-explained/).

## Serving from Files

When serving from zone files you use the *file* plugin. Let's start with the zone
`example.org.` and zonefile we want to serve from:

Create a file `example.org` with the following content:

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

Create a Corefile, `Corefile`, with:

~~~ txt
example.org {
    file example.org
    prometheus     # enable metrics
    errors         # show errors
    log            # enable query logs
}
~~~

Start CoreDNS on a non-standard port to check if everything is correct: `coredns -conf Corefile
-dns.port 1053` and send it a query with [dig](https://en.wikipedia.org/wiki/Dig_(command)):
~~~
% dig -p 1053 @localhost AAAA www.example.org +noall +answer

www.example.org.	3600	IN	AAAA	::1
~~~

As we've enabled query loggin with the [*log* plugin](/plugins/log) the query should be show up on
standard output as well:

~~~ txt
::1 - [24/Jul/2017:10:10:44 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 121 133.449µs
~~~

From here you can enable CoreDNS to run on port 53 and have it start from systemd (when on Linux),
see [the deployment repo](https://github.com/coredns/deployment) for example scripts.
Read more about the [*file*](/plugins/file/), [*metrics*](/plugins/metrics) and
[*errors*](/plugins/errors) plugin.

## CoreDNS as proxy

Another plugin is the [*proxy*](/plugins/proxy) plugin. We can for instance send DNS request to
Google over HTTPS. Create a Corefile with:

~~~ txt
. {
    proxy . 8.8.8.8:53 {
        protocol https_google
    }
    prometheus
    errors
    log
}
~~~

Start CoreDNS, just like above and send it a few queries. CoreDNS should logs those, in this case:
~~~
::1 - [24/Jul/2017:10:44:15 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 76 83.396955ms
::1 - [24/Jul/2017:10:44:17 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 76 14.030914ms
::1 - [24/Jul/2017:10:44:19 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 76 13.286384ms
~~~

If you look at the time each query took (in "ms") it's quite slow, ~83ms, 13ms. So
let's add some caching and
enable the [*caching*](/plugins/cache) plugin. Just add the word "cache" to the Corefile and
graceful reload CoreDNS: `kill -SIGUSR1 <pid_of_coredns>`. And query again:

~~~
::1 - [24/Jul/2017:11:33:54 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 76 43.469743ms
::1 - [24/Jul/2017:11:33:55 +0000] "AAAA IN www.example.org. udp 45 false 4096" NOERROR 73 133.073µs
~~~

First one is still "slow", but the subsequent query only takes 133 µs.

## Possible Errors

The [*health*](/plugins/health)'s documentation states "This plugin only needs to be enabled once",
which might lead you to think that this would be a valid Corefile:

~~~ txt
health

. {
    whoami
}
~~~
But this doesn't work and leads to the somewhat cryptic error:

    "Corefile:3 - Error during parsing: Unknown directive '.'".

What happens here? `health` is seen as zone and now the
parser expect to see directives (`cache`, `etcd`, etc.), but instead the next token is `.`, which
isn't a directive. The Corefile should be constructed as follows:

~~~ corefile
. {
    whoami
    health
}
~~~
That line in the *health*'s documentation means that once *health* is specified, it is global for
the entire CoreDNS process, even though you've only specified it for one server.

## Also See

There are [numerous other](/plugins) plugins that can be used with CoreDNS. And you can write [your
own](https://coredns.io/2016/12/19/writing-plugins-for-coredns/) plugin.

How [queries are processed](https://coredns.io/2017/06/08/how-queries-are-processed-in-coredns/) is
a deep dive into how CoreDNS handles DNS queries.
