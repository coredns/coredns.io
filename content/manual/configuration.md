# Configuration

There are various pieces that can be configured in CoreDNS. The first is determining which
plugins you want to compile into CoreDNS. The binaries we provide have all plugins, as listed in
[`plugin.cfg`](https://github.com/coredns/coredns/blob/master/plugin.cfg), compiled in.
Adding or removing is [easy](/2017/07/23/add-external-plugins/), but requires a recompile of CoreDNS.

Thus most users use the *Corefile* to configure CoreDNS. When CoreDNS starts, and the `-conf` flag is
not given, it will look for a file named `Corefile` in the current directory. That file consists
of one or more Server Blocks. Each Server Block lists one or more Plugins. Those Plugins may be
further configured with Directives.

The ordering of the Plugins in the Corefile *does not determine* the order of the plugin chain. The
order in which the the plugins are executed is determined by the ordering in `plugin.cfg`.

Comments in a Corefile are started with a `#`. The rest of the line is then considered a comment.

## Environment Variables

CoreDNS supports environment variables in its configuration.
They can be used anywhere in the Corefile. The syntax is `{$ENV_VAR}` (a more Windows-like syntax
`{%ENV_VAR%}` is also supported). CoreDNS substitutes the contents of the variable while parsing
the Corefile.

## Importing Other Files

See the [*import*](https://coredns.io/explugins/import) plugin. This plugin is a bit special in that
it may be used anywhere in the Corefile.

### Reusable Snippets

A special case of importing files is a *snippet*. A snippet is defined by naming a block with
a special syntax. The name has to be put in parentheses: `(name)`. After that, it can be included in
other parts of the configuration with the *import* plugin:

~~~ corefile
# define a snippet
(snip) {
    prometheus
    log
    errors
}

. {
    whoami
    import snip
}
~~~

## Server Blocks

Each Server Block starts with the zones the Server should be authoritative for. After the zone
name or a list of zone names (separated with spaces), a Server Block is opened with an opening brace.
A Server Block is closed with a closing brace. The following Server Block specifies a server that is
responsible for all zones below the root zone: `.`; basically, this server should handle every
possible query:

~~~ corefile
. {
    # Plugins defined here.
}
~~~

Server blocks can optionally specify a port number to listen on. This defaults to port 53 (the
standard port for DNS). Specifying a port is done by listing the port after the zone separated by
a colon. This Corefile instructs CoreDNS to create a Server that listens on port 1053:

~~~ corefile
.:1053 {
    # Plugins defined here.
}
~~~

> Note: if you explicitly define a listening port for a Server you *can't* overrule it with the
> `-dns.port` option.

Specifying a Server Block with a zone that is already assigned to a server *and* running it on the
same port is an error. This Corefile will generate an error on startup:

~~~ txt
.:1054 {

}

.:1054 {

}
~~~

Changing the second port number to 1055 makes these Server Blocks two different Servers.

### Specifying a Protocol

Currently CoreDNS accepts three different protocols: plain DNS, DNS over TLS and DNS over gRPC. You
can specify what a server should accept in the server configuration by prefixing a zone name with
a scheme.

* `dns://` for plain DNS (the default if no scheme is specified).
* `tls://` for DNS over TLS.
* `grpc://` for DNS over gRPC.

## Plugins

Each Server Block specifies a number of plugins that should be chained for this specific Server. In
its most simple form, you can add a Plugin by just using its name in a Server Block:

~~~ corefile
. {
    chaos
}
~~~

The *chaos* plugin makes CoreDNS answer queries in the CH class - this can be useful for identifying
a server. With the above configuration, CoreDNS will answer with its version when getting a request:

~~~ sh
$ dig @localhost -p 1053 CH version.bind TXT
...
;; ANSWER SECTION:
version.bind.		0	CH	TXT	"CoreDNS-1.0.5"
...
~~~

Most plugins allow more configuration with Directives. In the case of the [*chaos*](/plugins/chaos)
plugin we can specify a `VERSION` and `AUTHORS` as shown in its syntax:

> #### Syntax
>
> ```
> chaos [VERSION] [AUTHORS...]
> ```
>
> * **VERSION** is the version to return. Defaults to `CoreDNS-<version>`, if not set.
> * **AUTHORS** is what authors to return. No default.

So, this adds some Directives to the *chaos* plugin that will make CoreDNS will respond with
`CoreDNS-001` as its version:

~~~ corefile
. {
    chaos CoreDNS-001 info@coredns.io
}
~~~

Other plugins with more configuration options have a Plugin Block, which, just as a Server
Block, is enclosed in an opening and closing brace:

~~~ txt
. {
    plugin {
       # Plugin Block
    }
}
~~~

We can combine all this and have the following Corefile, which sets up 4 zones serving on two
different ports:

~~~ corefile
coredns.io:5300 {
    file db.coredns.io
}

example.io:53 {
    log
    errors
    file db.example.io
}

example.net:53 {
    file db.example.net
}

.:53 {
    kubernetes
    proxy . 8.8.8.8
    log
    errors
    cache
}
~~~

When parsed by CoreDNS, this will result in the following setup:

![CoreDNS: Zones, plugins and query routing](/images/CoreDNS-Corefile.png)

## External Plugins

External plugins are plugins that are not compiled into the default CoreDNS. You can easily enable
them, but you'll need to compile CoreDNS your self.

## Possible Errors

The [*health*](/plugins/health) plugin's documentation states "This plugin only needs to be enabled once",
which might lead you to think that this would be a valid Corefile:

~~~ txt
health

. {
    whoami
}
~~~
But this doesn't work and leads to a somewhat cryptic error:

~~~
"Corefile:3 - Error during parsing: Unknown directive '.'".
~~~

What happened here? `health` is seen as a zone (and the start of a Server Block). The parser expects to
see plugin names (`cache`, `etcd`, etc.), but instead the next token is `.`, which isn't a plugin.
The Corefile should be constructed as follows:

~~~ corefile
. {
    whoami
    health
}
~~~
That line in the *health* plugin's documentation means that once *health* is specified, it is global for
the entire CoreDNS process, even though you've only specified it for one server.
