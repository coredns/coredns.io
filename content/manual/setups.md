# Setups

Here you can find a bunch of configurations for CoreDNS. All setups are done assuming you are not the
root user and hence can't start listening on port 53. We will use port 1053 instead, using the
`-dns.port` flag. In every setup, the configuration file used is the CoreDNS' default, named `Corefile`.
This means we *don't need* to specify the configuration file with the `-conf` flag. In other words,
we start CoreDNS with `./coredns -dns.port=1053 -conf Corefile`, which can be abbreviated to
`./coredns -dns.port=1053`.

All DNS queries will be generated with the [`dig`](https://en.wikipedia.org/wiki/Dig_(command))
tool, the gold standard for debugging DNS. The *full* command line we use here is:
~~~ sh
$ dig -p 1053 @localhost +noall +answer <name> <type>
~~~
But we shorten it in the setups below, so `dig www.example.org A` is really
`dig -p 1053 @localhost +noall +answer www.example.org A`

## Authoritative Serving From Files

This setup uses the [*file*](/plugins/file) plugin. Note the external [*redis*](/explugins/redis) plugin
enables authoritative serving from a Redis Database. Let's continue with the setup using *file*.

The file we create here is a DNS zone file, and it can have any name (*file* plugin doesn't care). The
data we are putting in the file is for the zone `example.org.`.

In your current directory, create a file named `db.example.org` and put the following contents in
it:

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
The last two lines are defining a name `www.example.org.` with two addresses, 127.0.0.1 and (the
IPv6) ::1.

Next, create this minimal `Corefile` that handles queries for this domain and adds the
[*log*](/plugins/log) plugin to enable query logging:

~~~ corefile
example.org {
    file db.example.org
    log
}
~~~

Start CoreDNS and query it with `dig`:

~~~ sh
$ dig www.example.org AAAA

www.example.org.    3600    IN  AAAA    ::1
~~~

It works. Because of the *log* plugin, we should also see the query being logged:

~~~ txt
[INFO] [::1]:44390 - 63751 "AAAA IN www.example.org. udp 45 false 4096" NOERROR qr,aa,rd,ra 121 0.000106009s
~~~

The above logs show us the address CoreDNS replied from (`::1`) and the time and date it replied.
Furthermore, it logs the query type, the query class, the query name, the protocol used (`udp`), the
size in bytes of the incoming request, the DO bit state, and the advertised UDP buffer size. This is
data from the incoming query. `NOERROR` signals the start of the reply, which is the Response Code
sent back, followed by the set of flags on the reply: `qr,aa,rd,ra`, the size of the reply in bytes
(121), and the duration it took to get the reply.

## Forwarding

CoreDNS can be configured to forward traffic to a recursor with the [*forward*](/plugins/forward).
Here, we will use *forward* and focus on the most basic setup: forwarding to Google Public DNS
(8.8.8.8) and Quad9 DNS (9.9.9.9).

We don't need to create anything except for a Corefile with the configuration we want. In
this case, we want *all* queries hitting CoreDNS to be forward to either 8.8.8.8 or 9.9.9.9:

~~~ corefile
. {
    forward . 8.8.8.8 9.9.9.9
    log
}
~~~
Note that *forward* allows you to fine tune the names it will send upstream. Here, we
chose all names (`.`). For instance: `forward example.com 8.8.8.8 9.9.9.9` would only forward names
within the `example.com.` domain.

Start CoreDNS and test it with `dig`:
~~~ sh
$ dig www.example.org AAAA
www.example.org.	25837	IN	AAAA	2606:2800:220:1:248:1893:25c8:194
~~~

## Forwarding Domains To Different Upstreams

A common scenario you may encounter is that queries for `example.org` need to go to 8.8.8.8 and
the rest should be resolved via the name servers in `/etc/resolv.conf`. There are two ways that
could be implemented in a Corefile:

1. Use multiple declarations of [*forward*](/plugins/forward) plugin within a single Server Block:

~~~ txt
. {
    forward example.org 8.8.8.8
    forward . /etc/resolv.conf
    log
}
~~~

2. Use multiple Server Blocks, one for each of the domains you want to route on:

~~~ corefile
example.org {
    forward . 8.8.8.8
    log
}

. {
    forward . /etc/resolv.conf
    log
}
~~~

This leaves the domain routing to CoreDNS, which also handles special cases like DS queries. Having
two smaller Server Blocks instead of one has no negative effects except that your Corefile will be
slightly longer. Things like snippets and the [*import*](/plugins/import) will help there.

## Recursive Resolver

CoreDNS does not have a native (i.e. written in Go) recursive resolver, but there is an (external)
plugin that utilizes [libunbound](https://www.unbound.net/). For this setup to work, you first
have to recompile CoreDNS and [enable the *unbound*
plugin](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/). Super quick
primer here (you must have the CoreDNS [source](#source) installed):

* Add `unbound:github.com/coredns/unbound` to `plugin.cfg`.
* Do a `go generate`, followed by `make`.

Note: the *unbound* plugin needs cgo to be compiled, which also means the coredns binary is now
linked against libunbound and not a static binary anymore.

Assuming this worked, you can then enable *unbound* with the following Corefile:

~~~ corefile
. {
    unbound
    cache
    log
}
~~~
[*cache*](/plugins/cache) has been included, because the (internal) cache from *unbound* is disabled to allow the
cache's metrics to works just like normal.
