+++
date = "2016-10-19T19:09:32Z"
description = "CoreDNS-002 Release Notes"
release = "002"
tags = ["Release", "002", "Notes"]
title = "CoreDNS-002 Release"
author = "miek"
+++

CoreDNS-002 has been [released](https://github.com/coredns/coredns/releases)!

CoreDNS is a DNS server that chains middleware, where each middleware implements a DNS feature.

## What is New

* `-port` was renamed to `-dns.port` to avoid clashing with Caddy's `-port` (which was renamed to
  `http.port`).
* Lumberjack logger was removed, this means no built in log rotation; use an external tool for that.
* Brushed up GoDoc for all packages.
* Brushed up all READMEs to be more standard and look like manual page.
* Golint-ed and go vet-ed the code - these can now (somewhat) useful tools before submitting PRs.
* Add more tests and show test coverage on submitting/PRs.
* Various Corefile parsing bugs fixed, better syntax error detection.

## Middleware improvements:

* middleware/root: a root middleware, same usage as in [Caddy](https://caddyserver.com/docs/root).
  See middleware/root/README.md for its use in CoreDNS.
  This makes stanzas like this shorter:

    ~~~ txt
    .:53 {
        file /etc/coredns/zones/db.example.net example.net
        file /etc/coredns/zones/db.example.org example.org
        file /etc/coredns/zones/db.example.com example.com
    }
    ~~~

    Can be written as:

    ~~~ txt
    .:53 {
        root /etc/coredns/zones
        file db.example.net example.net
        file db.example.org example.org
        file db.example.com example.com
    }
    ~~~

* middleware/auto: similar to the *file* middleware, but automatically picks up new zones.
  The following Corefile will load all zones found under `/etc/coredns/org` and be authoritative
  for `.org.`:

    ~~~ txt
    .:53 {
        auto org {
            directory /etc/coredns/org
        }
    }
    ~~~
* middleware/file: handle wildcards better.
* middleware/kubernetes: TLS support for kubernetes and other improvements.
* middleware/cache: use an LRU cache to make it memory bounded. Added more option to have more
  control on what is cached and for how long. The cache stanza was extended:

    ~~~ txt
    cache {
        success CAPACITY [TTL]
        denial CAPACITY [TTL]
    }
    ~~~

  See middleware/cache/README.md for more details.

* middleware/dnssec: replaced go-cache with golang-lru in dnssec. Also adds a `cache_capacity`.
  option in dnssec middleware so that the capacity of the LRU cache could be specified in the config
  file.
* middleware/logging: allow a response classs to be specified on log on responses matching the name *and*
  the response class. For instance only log denials for example.com:

    ~~~ txt
    log example.com stdout {
        class denial
    }
    ~~~

* middleware/proxy: performance improvements.

# Contributors

The following people helped with getting this release done:

Chris O'Haver,
Manuel de Brito Fontes,
Miek Gieben,
Shawn Smith,
Silas Baronda,
Yong Tang,
Zhipeng Jiang.

If you want to help, please check out one of the [issues](https://github.com/coredns/coredns/issues/)
and start coding!

For documentation and help, see our [community page](https://coredns.io/community/).
