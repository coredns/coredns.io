+++
title = "finalize"
description = "*finalize* - resolves CNAMEs to their IP address."
weight = 10
tags = [  "plugin" , "finalize" ]
categories = [ "plugin", "external" ]
date = "2021-12-23T09:00:00+01:00"
repo = "https://github.com/tmeckel/coredns-finalizer"
home = "https://github.com/tmeckel/coredns-finalizer/blob/master/README.md"
+++

## Description

The plugin will try to resolve CNAMEs and only return the resulting A or AAAA
address. If no A or AAAA record can be resolved the original (first) answer will
be returned to the client.

Circular dependencies are detected and an error will be logged accordingly. In
that case the original (first) answer will be returned to the client as well.

## Syntax

```txt
finalize [max_depth MAX]
```

* `max_depth` **MAX** to limit the maximum calls to resolve a CNAME chain to the
    final A or AAAA record, a value `> 0` can be specified.

    If the maximum depth
    is reached and no A or AAAA record could be found, the the original (first)
    answer, containing the CNAME, will be returned to the client.

## Metrics

If monitoring is enabled (via the *prometheus* directive) the following metrics are exported:

* `coredns_finalize_request_count_total{server}` - query count to the *finalize* plugin.

* `coredns_finalize_circular_reference_count_total{server}` - count of detected circular references.

* `coredns_finalize_dangling_cname_count_total{server}` - count of CNAMEs that couldn't be resolved.

* `coredns_finalize_maxdepth_reached_count_total{server}` - count of incidents when max depth is reached while trying to resolve a CNAME.

* `coredns_finalize_maxdepth_upstream_error_count_total{server}` - count of upstream errors received.

* `coredns_finalize_request_duration_seconds{server}` - duration per CNAME resolve.

The `server` label indicated which server handled the request.

## Ready

This plugin reports readiness to the ready plugin. It will be immediately ready.

## Examples

In this configuration, we forward all queries to 9.9.9.9 and resolve CNAMEs.

```corefile
. {
  forward . 9.9.9.9
  finalize
}
```

In this configuration, we forward all queries to 9.9.9.9 and resolve CNAMEs with a maximum search depth of `1`:

```corefile
. {
  forward . 9.9.9.9
  finalize max_depth 1
}
```
