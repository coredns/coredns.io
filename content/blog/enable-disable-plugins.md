+++
title = "Compile Time Enabling or Disabling Plugins"
description = "Enable or Disable plugins when compiling CoreDNS."
tags = ["Documentation"]
draft = false
date = "2017-07-25T16:07:39+01:00"
author = "miek"
+++

CoreDNS' [plugins](/plugins) (or [external plugins](/explugins)) can be enabled or
disabled on the fly by specifying (or not specifying) it in the
[Corefile](/2017/07/23/corefile-explained/). But you can also compile CoreDNS with only the
plugin you *need* and leave the rest completely out.


There are two ways to achieve that. It could be done via compile-time configuration file
with CoreDNS code base update. It also could be achieved without modifying CoreDNS code.

## Build with compile-time configuration file

The with compile-time configuration file,
[`plugin.cfg`](https://github.com/coredns/coredns/blob/master/plugin.cfg) is all you need
to update. It looks like this:

~~~
...
whoami:whoami
erratic:erratic
startup:github.com/mholt/caddy/startupshutdown
...
~~~

The ordering of the plugins is specified by how the are ordered in this file. Each line consists of
a **name** and a **repository**. Just add or remove your plugin in this file.

Then do a `go get <plugin-repo-path>` if you need to get the external plugin's source code. And then
just compile CoreDNS with `go generate` and a `go build`. You can then check if CoreDNS has the new
plugin with `coredns -plugins`.

## Build with external golang source code

Alternatively, you could assembly plugins from different places through an external golang program.
It looks like this:

~~~
package main

import (
        _ "github.com/coredns/example"

        "github.com/coredns/coredns/coremain"
        "github.com/coredns/coredns/core/dnsserver"
)

var directives = []string{
        "example",
        ...
        ...
        "whoami",
        "startup",
        "shutdown",
}

func init() {
        dnsserver.Directives = directives
}

func main() {
        coremain.Run()
}
~~~

In the above sample code, the external plugin `example` has been imported with:
~~~
        _ "github.com/coredns/example"
~~~

The directives should also be updated through:
~~~
        dnsserver.Directives = directives
~~~

The ordering of the plugins is specified by how the arey ordered in the slice `directives`.

Then you can just compile CoreDNS with `go build` to have the binary generated with the
plugins you selected.
