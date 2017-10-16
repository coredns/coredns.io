+++
date = "2016-10-13T08:50:13Z"
description = "Which plugin will handle a query?"
tags = ["Corefile", "Query", "Routing", "Documentation"]
title = "Query Routing"
author = "miek"
+++

Quiz time, in the following Corefile:

~~~ txt
. {
  proxy . 8.8.8.8:53
  file db.example.com
}
~~~

Will a query for `www.google.com` be handled by the *proxy* or the *file* plugin? Answer below.

What does this Corefile actually say? It specifies that queries for root (`.`) and *everything*
below it (so for all domain names) we should enter this stanza.

Next *all* queries should be forwarded to 8.8.8.8:53.

Then because the *file* plugin **does not** specify what zones should be answered from the
`db.example.com` file, the toplevel one applies, which is root (`.`)

So we are left with a situation where both plugins will be called for the same names (which can
be perfectly valid for plugin that calls other chained-in plugin).

But *proxy* **will not** call *file* because the query will be answered and done with after
the plugin exists - the same is true for the opposite direction.

To look what into what happens here we have to look the [plugins
ordering](https://github.com/coredns/coredns/blob/master/plugins.cfg):

~~~
...
dnssec:dnssec
file:file
etcd:etcd
proxy:proxy
...
~~~

And we see that *file* is first and *proxy* comes somewhat later. This means that in the example
above *all* queries are routed to the *file* plugin. It will happily answer those with SERVFAIL,
because it probably can't find `www.google.com` in a file that will mostly have `*.example.com`
names in it.

In order to fix this, we should either have to separate stanza or specify the origin(s) for the
*file* plugin:

~~~ txt
. {
  proxy . 8.8.8.8:53
  file db.example.com example.com
}
~~~

To preempt a feature request: Yes, it would be nice of CoreDNS can detect and warn about this (it
does not do this now).
