# Writing Plugins

As mentioned before in this manual, plugins are the thing that make CoreDNS tick. We've seen
a bunch of configuration in the [previous section](#setups), but how can you write your own plugin?

See [Writing Plugins for CoreDNS](/2016/12/19/writing-plugins-for-coredns/) for an older post on
this subject. The [plugin.md](https://github.com/coredns/coredns/blob/master/plugin.md) documented
in CoreDNS's source also has some background and talks about styling the README.md.

The canonical example plugin is the *example* plugin. Its [github
repository](https://github.com/coredns/example) shows the most minimal code (with tests!) that is
needed to create a plugin.

It has:

1. `setup.go` and `setup_test.go`, which implement the parsing of configuration from the Corefile.
   The (usually named) `setup` function is called whenever the Corefile parser see the plugin's
   name; in this case, "example".
2. `example.go` (usually named `<plugin_name>.go`), which contains logic for handling the query, and
   `example_test.go`, which has basic units tests to check if the plugin works.
3. The `README.md` that documents in a Unix manual style how this plugin can be configured.
4. A LICENSE file. For inclusion in CoreDNS, this needs to have an APL like license.

The code also has extensive comments; feel free to fork it and base your plugin off of it.

## How Plugins Are Called

When CoreDNS wants to use a plugin it calls the method `ServeDNS`. `ServeDNS` has three parameters:

* a `context.Context`;
* a `dns.ResponseWriter` that is basically the client's connection;
* a `*dns.Msg` that is the request from the client.

`ServeDNS` returns two values: a (response) code and an error. The error is logged when the
[*errors*](/plugins/errors) is used in this server.

The code tells CoreDNS if a *reply has been written by the plugin chain or not*. In the latter case,
CoreDNS will take care of that. For the code's values, we *reuse* the DNS return codes (rcodes) from
the [dns](github.com/miekg/dns) package.

CoreDNS treats:

* SERVFAIL (dns.RcodeServerFailure)
* REFUSED (dns.RcodeRefused)
* FORMERR (dns.RcodeFormatError)
* NOTIMP (dns.RcodeNotImplemented)

as special and will then assume *nothing* has been written to the client. In all other cases, it
assumes something has been written to the client (by the plugin).

See [this post](https://blog.coredns.io/2017/03/01/how-to-add-plugins-to-coredns/) on how to compile
CoreDNS with your plugin.

## Logging From a Plugin

If your plugin needs to output a log line, you should use the `log` package. CoreDNS does not
implement log levels. The standard way of outputting is: `log.Printf("[LEVEL] ...")`, and LEVEL
can be: `INFO`, `WARNING` or `ERROR`.

In general, logging should be left to the higher layers when returning an error. However, if there is
a reason to consume the error but still notify the user, then logging in the plugin can be acceptable.

## Metrics

When exporting metrics, the *Namespace* should be `plugin.Namespace` (="coredns"), and the
*Subsystem* should be the name of the plugin. The README.md for the plugin should then also contain
a *Metrics* section detailing the metrics. If the plugin supports dynamic [health](/plugin/health)
reporting, it should also have a *Health* section detailing some of its inner workings.

## Documentation

Each plugin should have a README.md explaining what the plugin does and how it is configured. The
file should have the following layout:

* Title: use the plugin's name
* Subsection titled: "Named"
    with `<plugin name> - <one line description>`, i.e. NAME DASH DESCRIPTION
* Subsection titled: "Description" with a longer description and all the options the plugin supports.
* Subsection titled: "Syntax" detailing syntax and supported directives.
* Subsection titled: "Examples".
* Optional Subsection titled: "See Also", that references external documentation, like IETF RFCs.
* Optional Subsection titled: "Bugs" that lists things that do not work yet.

More sections are, of course, possible.

### Style

We use the Unix manual page style:

* The name of the plugin in the running text should be italic: `*plugin*`.
* All CAPITAL user supplied arguments in the running text reference use strong text: `**EXAMPLE**`.
* Optional text is in block quotes: `[optional]`.
* Use three dots to indicate multiple options are allowed: `arg...`.
* Item used literal: `literal`.

### Example Domain Names

Please be sure to use `example.org` or `example.net` in any examples and tests you provide. These
are the standard domain names created for this purpose. If you don't, there is a chance your fantasy
domain name is registered by someone and will actually serve web content (which you may like or not).

## Fallthrough

In a perfect world, the following would be true for plugins: "Either you are responsible for a zone or
not". If the answer is "not", the plugin should call the next plugin in the chain. If "yes" it
should handle *all* names that fall in this zone and the names below - i.e. it should handle the
entire domain and all sub domains.

TODO(miek): ref to "Query Is Proccessed with Fallthrough"

~~~ txt
. {
    file example.org db.example
}
~~~

In this example the *file* plugin is handling all names below (and including) `example.org`. If
a query comes in that is not a subdomain (or equal to) `example.org` the next plugin is called.

Now, the world isn't perfect, and there are good reasons to "fallthrough" to the next middleware,
meaning a plugin is only responsible for a *subset* of names within the zone. The first of these
to appear was the *reverse* plugin that synthesizes PTR and A/AAAA responses (useful with IPv6).

The nature of the *reverse* plugin is such that it only deals with A/AAAA and PTR, and then only
for a subset of the names. Ideally, you would want to layer *reverse* **in front of** another
plugin such as *file* or *auto* (or even *proxy*). This means *reverse* handles some special
reverse cases and **all other** requests are handled by the backing plugin. This is exactly what
"fallthrough" does. To keep things explicit we've opted that plugins implementing such behavior
should implement a `fallthrough` keyword.

The `fallthrough` directive should optionally accept a list of zones. Only queries for records
in one of those zones should be allowed to fallthrough.

## Qualifying for main repo

Plugins for CoreDNS can live out-of-tree. `plugin.cfg` defaults to CoreDNS's repo, but external
repos work fine. So when do we consider the inclusion of a new plugin in the main repo?

* The plugin authors should be willing to maintain the plugin, i.e. your GitHub handle will be
  listed in its `OWNERS` file.
* The plugin should be useful for other people. "Useful" is a subjective term, but it should
  bring something new to CoreDNS.
* It should be sufficiently different from other plugins to warrant inclusion.
* Current Internet standards need be supported: IPv4 and IPv6, so A and AAAA records should be
  handled (if your plugin is in the business of dealing with address records, that is).
* It must have tests.
* It must have a README.md for documentation.
