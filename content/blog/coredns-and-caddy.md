+++
date = "2016-09-29T06:18:40Z"
description = "CoreDNS can now be downloaded from Caddy's download page."
tags = ["Caddy", "Bugs", "Questions"]
title = "CoreDNS and Caddy"
author = "miek"
+++
[Caddy 0.9.3](https://forum.caddyserver.com/t/caddy-0-9-3-released/725) is released.  On it's
[download page](https://caddyserver.com/download) you can now select the "DNS plugin" to be added
to Caddy! This is really nice and a culmination of all the work that has been put in to make this
happen.

Note that if you select this option you get a binary that is *both* a DNS and webserver, during
startup you can select between the two with `-type=dns|http` flag.

The CoreDNS developers will still provide their own (DNS only) binaries over [at
github](https://github.com/coredns/coredns/releases).

Note that we *also* have a CoreDNS category on the [Caddy
forum](https://forum.caddyserver.com/c/coredns) where you can ask questions. Bugs and feature
requests are probably better directed to our [github](https://github.com/coredns/coredns).
