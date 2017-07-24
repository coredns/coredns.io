+++
date = "2017-07-23T21:28:13Z"
description = "How to add a link to external middleware?"
tags = ["Middleware", "External", "Add"]
title = "Add External Middleware"
author = "miek"
+++

If you want to have your external middleware listed here send an email to the coredns-discuss email
list with the details or create [a pull request](https;//github.com/coredns/coredns.io).

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
repo = "link-to-your-repo"
home = "link-to-your-homepage"
+++

The *<middleware name>* middleware is a ...

<rest of your README.md>
~~~

See [example.md for an example]
(https://github.com/coredns/new.coredns.io/blob/master/content/exmiddleware/example.md) on how to do
this.
