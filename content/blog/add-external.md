+++
date = "2017-07-23T21:28:13Z"
description = "How to add a link to external middleware?"
tags = ["Middleware", "External", "Add"]
title = "Add External Middleware"
author = "miek"
+++

If you want to have your external middleware listed here send an email to the [coredns-discuss email
list](/community) with the details or create [a pull
request](https://github.com/coredns/coredns.io).

In that email or pull request you'll need to add a file to `content/blog/exmiddleware/` that looks
like this:

~~~ txt
+++
title = "<middleware name>"
description = "*<middleware name>* is a ..."
weight = 10
tags = [  "middleware" , "<middleware name>" ]
categories = [ "middleware", "external" ]
date = "2017-07-22T12:37:19+01:00"
repo = "https://link-to-your-middleware-repo"
home = "https://link-to-your-homepage-or-readme"
+++

The *<middleware name>* middleware is a ...

<rest of your README.md>
~~~

See [shutdown.md for an example]
(https://github.com/coredns/new.coredns.io/blob/master/content/exmiddleware/shutdown.md) on how to do
this.

Note that **description** needs to be a full sentence, and that **repo** must be a Go-gettable link
that can be put in `middleware.cfg`. Also document in your README.md what would be suitable position
for your middleware in `middleware.cfg` (beginning, middleware, end or give it a specific number).
