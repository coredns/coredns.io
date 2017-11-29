+++
date = "2016-11-26T17:22:44Z"
description = "Using Google's dns.google.com with CoreDNS."
tags = ["Encryption", "DNS", "Google", "Documentation"]
title = "DNS over HTTPS"
author = "miek"
+++

Since almost a year Google has a DNS service that can be queried over HTTPS:
<https://dns.google.com>. This means your queries are encrypted and can only be seen by you (and
Google(!)). Seeing all the press about the
[UK's snooper's charter](https://www.theguardian.com/technology/askjack/2016/nov/24/how-can-i-protect-myself-from-government-snoopers)
I though I should implement this as a plugin in CoreDNS.

I'm (obviously) going to use this myself; which is perfect as it protects me and it allows me to
dog food CoreDNS as a DNS proxy in my home network.

A note worthy other implementation is "dingo": <https://github.com/pforemski/dingo>.

Also note that this a *different* protocol than "DNS over TLS" which has similar goals and is being
standardized by the IETF.

> Currently you'll need to compile CoreDNS from source to play with this or wait until CoreDNS-004
> is released.

The configuration on the CoreDNS side is pretty straight forward. The following Corefile is all
you'll need:

~~~ corefile
. {
    proxy . 8.8.8.8 {
        protocol https_google
    }
    cache
    log
    errors
}
~~~

Next start CoreDNS, and query it.

~~~ sh
% ./coredns
.:53
2016/11/26 17:11:07 [INFO] CoreDNS-003
CoreDNS-003
::1 - [26/Nov/2016:17:13:10 +0000] "MX IN miek.nl. udp false 4096" NOERROR 246 149.791162ms
::1 - [26/Nov/2016:17:13:11 +0000] "MX IN miek.nl. udp false 4096" NOERROR 170 156.432µs
~~~

The only unencrypted DNS used is from your laptop/phone/computer to CoreDNS, the rest is encrypted.

By default, `dns.google.com` will be re-resolved every 30 seconds using 8.8.8.8 and 8.8.4.4 (you can
override these defaults). This is the only query not encrypted, but this will probably lead to
a very boring browser history.

Next, *I* need to configure a Raspberry Pi and install CoreDNS on it. And as with all CoreDNS
developements [feedback is welcome](https://github.com/coredns/coredns/issues).
