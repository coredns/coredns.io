+++
title = "warnlist"
description = "*warnlist* emits logs and Prometheus metrics when a listed domain is requested."
weight = 10
tags = [  "plugin" , "warnlist" ]
categories = [ "plugin", "external" ]
date = "2021-06-03T16:00:00+01:00"
repo = "https://github.com/giantswarm/coredns-warnlist-plugin"
home = "https://github.com/giantswarm/coredns-warnlist-plugin/blob/main/README.md"
+++

## Description

The *warnlist* plugin accepts a list of malicious or otherwise undesirable domains and emits a log entry and Prometheus metrics when a domain (or subdomain) is requested.

Prohibited domains can be loaded from a local file or a URL and can be automatically reloaded after a specified period.

*warnlist* can be thought of as a non-blocking blacklist/blocklist/denylist/badlist. When used with a curated data source, the plugin can surface simplistic low-noise alerts without the need to ship and inspect DNS logs.

Refer to the project README for more info.

An unofficial `coredns` image with this plugin already compiled is hosted by [Giant Swarm](giantswarm.io) on [Quay](https://quay.io/repository/giantswarm/coredns-warnlist-plugin) and [Docker Hub](https://hub.docker.com/r/giantswarm/coredns-warnlist-plugin), as `[quay.io/]giantswarm/coredns-warnlist-plugin`.

## Syntax

```text
warnlist {
    <source type> <source path> <file format>
    reload <reload period>
    match_subdomains <true | false>
}
```

The `warnlist` plugin accepts the following arguments:

- `<source type>`: Type of the domain list. Either `url` or `file`.
- `<source path>`: Where to load the list from. Either a URL or file path.
- `<file format>`: Format of the file to expect. Either `hostfile` or `text`.
- `<reload period>`: (Optional) Go Duration after which the list will be regenerated*.
- `<match subdomains>`: (Optional) If `true` (default), the plugin will also check and match subdomains of those explicitly listed. Either `true` or `false`.

\* A jitter of +/- 30% is automatically added. When automatically reloading from a URL, please be friendly to the service hosting the file.

## Example

Sample `Corefile` using a URL data source, reloading every ~60 minutes:

```text
. {
    log
    warnlist {
        url https://urlhaus.abuse.ch/downloads/hostfile/ hostfile
        reload 60m
    }
    prometheus
    forward . /etc/resolv.conf
}
```

## Metrics

If the `prometheus` plugin is also enabled, this plugin emits the following metrics:

- `warnlist_hits_total{server, requestor, domain}` - counts the number of warnlisted domains requested. The host and domain are included as labels.
- `warnlist_failed_reloads_count{server}` - counts the number of times the plugin has failed to reload.
- `warnlist_cache_check_duration_seconds{server}` - summary for determining the average time it takes to check the warnlist.
- `warnlist_warnlisted_items_count{server}` - current number of domains stored in the warnlist.
