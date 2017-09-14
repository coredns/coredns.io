+++
date = "2017-07-23T21:28:13Z"
description = "How to add a link to external plugin?"
tags = ["Plugin", "External", "Add"]
title = "Add External Plugins"
author = "miek"
+++

If you want to have your external plugin listed here send an email to the [coredns-discuss email
list](/community) with the details or create [a pull
request](https://github.com/coredns/coredns.io).

In that email or pull request you'll need to add a file to `content/blog/explugins/` that looks
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

The *<plugin name>* plugin is a ...

## Syntax

## Example

## How to Enable

Follow [these](/2017/07/25/compile-time-enabling-or-disabling-plugins/) steps.
~~~

See [shutdown.md for an example]
(https://github.com/coredns/coredns.io/blob/master/content/explugins/shutdown.md) on how to do
this.

Note that **description** needs to be a full sentence, and that **repo** must be a Go-gettable link
that can be put in `plugins.cfg`. Also document in your README.md what would be suitable position
for your plugin in `plugins.cfg` (beginning, middle, end or give it a specific number).
