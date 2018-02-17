# Configuration

There are to various pieces that can be configured in CoreDNS. The first is determining which
plugins you want to compile into CoreDNS. The binaries we provide have all plugins as listed in
`plugin.cfg` compiled in. Adding or removing is [easy](/link/to/howto), but shouldn't normally be
done by end users (or shouldn't needed).

Thus most users use the *Corefile* to configure CoreDNS. When CoreDNS starts, and the `-conf` flag is
not given it will look for a file named `Corefile` in the current directory. That files consists out
of one or more Server stanzas. Each Server stanza lists one or more Plugins. Those Plugins may be
further configured with Directives. The ordering of the Plugin in the Corefile *does not determine*
the order of the plugin chain.

As said (/link) the plugin chain ordering is fixed and determined via plugin.cfg during the
compilation phase.

Comments in a Corefile are started with a `#`. The rest of the line is then considered a comment.

## Environment Variables

CoreDNS supports [environment variables](no-caddy-ref).
They can be used anywhere in the Corefile. The syntax is `{$ENV_VAR}` (a more Windows like syntax
`{%ENV_VAR%}` is also supported). CoreDNST substitutes the contents of the variable while parsing
the Corefile.

## Importing other files

See the [*import*](https://coredns.io/explugins/import) plugin. This plugin is a bit special in that
it may be used anywhere in the Corefile.

## Server Blocks

Each Server Block starts with the zones this Server should be authoritative for. After this zone
name or list of zone names (separated with spaces) a Server Block is opened with an opening brace.
(Sometimes a Server Block is referred to as a Server Stanza). A Server Block is closed with
a closing brace. The following Server Block specifies a server that is responsible for all zones
below the root zone: `.`, basically this server should handle every possible query:

~~~ corefile
. {
    # Plugins defined here.
}
~~~

Server blocks can optionally specify the port number to listen on. This defaults to port 53 (the
standard port for DNS). Specifying the port is done by listing after the zone separated with
a colon. This Corefile instructs CoreDNS to listen on port 1053.

~~~ corefile
.:1053 {
    # Plugins defined here.
}
~~~

> Note: if you explicitly define a listening port for a Server you *can't* overrule it with the
> `-dns.port` option.

Specifying a Server Block with a zone that is already assigned to a server *and* running it on the
same port is an error:

~~~ corefile
.:1054 {

}

.:1054 {

}
~~~

Will generate an error on startup. Changing the second port number to 1055 makes these Server Blocks
two different Servers.

## Specifying a Protocol

Currently CoreDNS accepts three different protocols: plain DNS, DNS over TLS and DNS over gRPC, you
can specify what a server should accept if the configuration, by prefixing a zone name with
a scheme, use:

* `dns://` for plain DNS (the default is no scheme is specified).
* `tls://` for DNS over TLS.
* `grpc://` for DNS over gRPC.


## Plugins

Each Server Block specifies a bunch of plugins (plugins.md)

## External Plugin

External plugins are plugins that are not compiled into the default CoreDNS. You can easily enable
them, but you'll need to compile CoreDNS your self.
