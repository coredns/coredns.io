+++
title = "Compile Time Enabling or Disabling Middleware"
description = "Enable or Disable middleware when compiling CoreDNS"
tags = ["Documentation"]
draft = false
date = "2017-07-25T16:07:39+01:00"
author = "miek"
+++

CoreDNS' [middleware](/middleware) (or [external middleware](/exmiddleware)) can be enabled or
disabled on the fly by specifying (or not specifying) it in the
[Corefile](/2017/07/23/corefile-explained/). But you can also compile CoreDNS with only the
middleware you *need* and leave the rest completely out.

All this is done via one compile-time configuration file,
[`middleware.cfg`](https://github.com/coredns/coredns/blob/master/middleware.cfg). It looks like this:

~~~
...
230:whoami:whoami
240:erratic:erratic
500:startup:github.com/mholt/caddy/startupshutdown
...
~~~

The number specifies the ordering of the middleware (they are called in this order - *if* enabled - by
CoreDNS). Then a **name** and a **repository**. Just add or remove your middleware in this file.

Then do a `go get` if you need to get the external middleware's source code. And then just compile
CoreDNS with `go generate` and a `go build`.
