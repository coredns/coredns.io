# Plugins

When CoreDNS has been started and has parsed the configuration it runs Servers.
Each Server is defined by the zones it serves and on what port. Each Server has
it's own Plugin Chain.

When a query is being processed by CoreDNS the following steps are performed:

1. If there are multiple Servers configured that listen on this port it will check which one
   has the most specific zone for this query (longest suffix match). E.g. if there are two Servers,
   one for `example.org` and one for `a.example.org`, and the query is for `www.b.example.org`, it
   will be routed to the latter.
2. Once a Server has been found, it will be routed through the Plugin Chain that is configured for
   this server. This always happens in the same order. That (static) ordering is defined in
   [`plugin.cfg`](https://github.com).
3. Each plugin will inspect the query and determine if it should process it (some plugins
   allow you to filter further on the query name or other attributes).
   A couple of things can now happen:

   1. The query is processed by this plugin.
   2. The query is *not* processed by this plugin
   3. The query is processed by this plugin, but half way through it decides it still want
      to call the next plugin in the chain. We call this *fallthrough* after the keyword that
      enables it.
   4. The query is processed by this plugin, a "hint" is added and the next plugin is called. This
      hint provides a way to "see" the (eventual) response and act up on that.

Processing a query means a Plugin will respond to the client with a reply.

Note that a plugin is free to deviate from the above list as it wishes. Currently all plugins that
come with CoreDNS fall into one of these four groups though.

## Query Is Processed

The plugin processes the query. It looks up (or generates, or whatever the plugin author decided
this plugin does) a response and sends it back to the client. The query processing stops here, no
next plugin is called. A (simple) plugin that works like this is [*whoami*](/plugins/whoami).

## Query Is Not Processed

If the plugin decides it will not process a query, it simply calls the next plugin in the chain.
When the last plugin in the chain also fails to process a query, CoreDNS will return SERVFAIL back
to the client.

## Query Is Proccessed with Fallthrough

In this situation a plugin handles the query, but the reply it got from its backend (i.e. maybe it
got NXDOMAIN), is such that it wants other plugins in the chain to a look as well. If *fallthrough*
is provided (and enabled!) the next plugin is called. A plugin that works like this is the
[*hosts*](/plugins/hosts) plugin.
First a lookup in the host table (`/etc/hosts`) is attempted, if it finds an answer it returns that.
If not, it will *fallthrough* to the next one in the hopes other plugins may return something to the
client.

## Query is Processes with a Hint

A plugin of this kind will process a query, and will *always* call the next plugin. But it provides
a hint that allows it to see the response that will be written to the client. A plugin that does
this is [*metrics*](/plugin/metrics). It times the duration ...

## Anatomy of Plugins

A plugin consists out of a Setup, Registration and Handler part.

The Setup parses the configuration and the Plugin's Directives (those should be documented in the
plugin's README).

The Handler is the code that processes the query and implement all the logic.

The Registration part registers the plugin in CoreDNS - this happens when CoreDNS is compiled. All
of the registred plugin can be used by a Server. Which plugin are configured in each Server happens
run time and is done in CoreDNS's configuration file. The Corefile.

## Plugin Documenation

Each plugin has it's own README detailing how it can be configured. This README includes example and
other bit a user should be aware of.
