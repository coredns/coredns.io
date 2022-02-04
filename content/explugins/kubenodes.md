+++
title = "kubenodes"
description = "*kubenodes* - creates records for Kubernetes nodes."
weight = 10
tags = [  "plugin" , "kubenodes" ]
categories = [ "plugin", "external" ]
date = "2021-12-17T00:00:00+00:00"
repo = "https://github.com/infobloxopen/kubenodes"
home = "https://github.com/infobloxopen/kubenodes/blob/main/README.md"
+++

## Description

*kubenodes* watches the Kubernetes API and synthesizes A, AAAA, and PTR records for Node addresses.

This plugin requires ...
* the [_kubeapi_ plugin](http://github.com/coredns/kubeapi) to create a connection
  to the Kubernetes API.
* list/watch permission to the Nodes API.

This plugin can only be used once per Server Block.

## Syntax

```
kubenodes [ZONES...] {
    external
    ttl TTL
    fallthrough [ZONES...]
}
```
* `external` will build records using Nodes' external addresses.  If omitted, *kubenodes* will build records using
  Nodes' internal addresses.
* `ttl` allows you to set a custom TTL for responses. The default is 5 seconds.  The minimum TTL allowed is
  0 seconds, and the maximum is capped at 3600 seconds. Setting TTL to 0 will prevent records from being cached.
  All endpoint queries and headless service queries will result in an NXDOMAIN.
* `fallthrough` **[ZONES...]** If a query for a record in the zones for which the plugin is authoritative
  results in NXDOMAIN, normally that is what the response will be. However, if you specify this option,
  the query will instead be passed on down the plugin chain, which can include another plugin to handle
  the query. If **[ZONES...]** is omitted, then fallthrough happens for all zones for which the plugin
  is authoritative. If specific zones are listed (for example `in-addr.arpa` and `ip6.arpa`), then only
  queries for those zones will be subject to fallthrough.

## External Plugin

To use this plugin, compile CoreDNS with this plugin added to the `plugin.cfg`.  It should be positioned before
the _kubernetes_ plugin if _kubenode_ is using the same zone or a superzone of _kubernetes_.  This plugin also requires
the _kubeapi_ plugin, which should be added to the end of `plugin.cfg`.

## Ready

This plugin reports that it is ready to the _ready_ plugin once it has received the complete list of Nodes
from the Kubernetes API.

## Examples

Use Nodes' internal addresses to answer forward and reverse lookups in the zone `node.cluster.local.`.
Fallthrough to the next plugin for reverse lookups that don't match any Nodes' internal IP addresses.

```
kubeapi
kubenodes node.cluster.local in-addr.arpa ip6.arpa {
  fallthrough in-addr.arpa ip6.arpa
}
```

Use Nodes' external addresses to answer forward and reverse lookups in the zone `example.`. Fallthrough
to the next plugin for reverse lookups that don't match any Nodes' external IP addresses.

```
kubeapi
kubenodes example in-addr.arpa ip6.arpa {
  external
  fallthrough in-addr.arpa ip6.arpa
}
```

