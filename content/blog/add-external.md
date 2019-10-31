+++
date = "2017-07-23T21:28:13Z"
description = "How to add a link to external plugin?"
tags = ["Plugin", "External", "Add"]
title = "Add External Plugins"
author = "miek"
+++

If you want to have your external plugin listed create [a pull request](https://github.com/coredns/coredns.io).

In that pull request you'll need to add a file to `content/explugins/` that looks
like this:

~~~ txt
+++
title = "<plugin name>"
description = "*<plugin name>* is a ..."
weight = 10
tags = [  "plugin" , "<plugin name>" ]
categories = [ "plugin", "external" ]
date = "2017-07-22T12:37:19+01:00"
repo = "https://link-to-your-plugin-repo"
home = "https://link-to-your-homepage-or-readme"
+++

## Description

The *<plugin name>* plugin is a ...

## Syntax

## Example
~~~

See [example.md for an example]
(https://raw.githubusercontent.com/coredns/coredns.io/master/content/explugins/example.md)
on how to do this.

Note that **description** needs to be a full sentence, and that **repo** must be a Go-gettable link
that can be put in `plugins.cfg`.

## Go Modules

Note we found the Go modules can interact badly with how external plugins are compiled into CoreDNS.
Various external plugins have *removed* the `go.mod` and `go.sum` files entirely to work around
issues, see for instance [unbound](https://github.com/coredns/unbound/).
