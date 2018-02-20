+++
date = "2017-07-23T21:28:13Z"
description = "How to add a link to external plugin?"
tags = ["Plugin", "External", "Add"]
title = "Add External Plugins"
author = "miek"
+++

If you want to have your external plugin listed create [a pull request](https://github.com/coredns/coredns.io).

In that pull request you'll need to add a file to `content/blog/explugins/` that looks
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

See [shutdown.md for an example]
(https://raw.githubusercontent.com/coredns/coredns.io/master/content/explugins/shutdown.md)
on how to do this.

Note that **description** needs to be a full sentence, and that **repo** must be a Go-gettable link
that can be put in `plugins.cfg`.
