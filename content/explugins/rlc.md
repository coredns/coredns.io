+++
title = "rlc"
description = "*rlc* - enrich reverse lookups with names observed in forward DNS answers."
weight = 10
tags = [ "plugin", "rlc" ]
categories = [ "plugin", "external" ]
date = "2026-10-29T00:00:00+00:00"
repo = "https://github.com/houyuwushang/coredns-rlc"
home = "https://github.com/houyuwushang/coredns-rlc/blob/main/README.md"
+++

## Description

*rlc* learns IP/name associations from successful A and AAAA answers, including
the original question and its CNAME chain, and uses them to answer later PTR
queries. It is intended for log and monitoring enrichment when ordinary reverse
DNS is missing or unhelpful.

These are observed names, not authoritative PTR records or proof of which host
a connection accessed. Several names can share an IP, and answers depend on
which lookups the plugin has observed. Results are bounded and sorted by name;
the first name has no special meaning.

On a hit, *rlc* answers before downstream plugins, including any real PTR answer
they would provide. On a miss, it passes the query down the chain. DNSSEC DO and
TSIG requests pass through unchanged, and synthetic answers do not set AA or AD.
Retention is independent of the original DNS TTL, and PTR queries do not renew it.

The plugin supports optional TCP DNS peer lookups, DNS-based peer discovery for
headless Kubernetes Services, and dnstap response ingestion. Peer sharing is
best effort; the optional listeners should be restricted to trusted networks.

## Installation

Add `rlc:github.com/houyuwushang/coredns-rlc` immediately before `cache:cache` in
`plugin.cfg`, then rebuild CoreDNS. This order lets *rlc* observe cached forward
answers and respond even if the regular cache holds an older negative PTR.

The [plugin README](https://github.com/houyuwushang/coredns-rlc/blob/main/README.md)
contains complete build instructions, configuration options, limits, peer and
dnstap examples, and validation commands.

## Examples

```corefile
.:1053 {
    rlc {
        ttl 30m
        capacity 50000
        max_names 16
    }
    cache 30
    forward . /etc/resolv.conf
}
```

This keeps up to 50,000 IP/name pairs for 30 minutes, with at most 16 names per
IP. The default maximum answer TTL is 5 seconds, bounded by remaining retention.
Choose an explicit upstream if `/etc/resolv.conf` points back to this server.

## See Also

This independently maintained plugin continues the use case proposed by
Carsten Zeumer in [CoreDNS #6459](https://github.com/coredns/coredns/pull/6459),
following the discussion's external-plugin direction.
