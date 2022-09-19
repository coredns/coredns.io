+++
title = "acl"
description = "*acl* enforces access control policies on source ip and prevents unauthorized access to DNS servers."
weight = 1
tags = ["plugin", "acl"]
categories = ["plugin"]
date = "2022-09-19T13:19:15.8771589"
+++

## Description

With `acl` enabled, users are able to block or filter suspicious DNS queries by configuring IP filter rule sets, i.e. allowing authorized queries or blocking unauthorized queries.


When evaluating the rule sets, _acl_ uses the source IP of the TCP/UDP headers of the DNS query received by CoreDNS.
This source IP will be different than the IP of the client originating the request in cases where the source IP of the request is changed in transit.  For example:
* if the request passes though an intermediate forwarding DNS server or recursive DNS server before reaching CoreDNS
* if the request traverses a Source NAT before reaching CoreDNS

This plugin can be used multiple times per Server Block.

## Syntax

```
acl [ZONES...] {
    ACTION [type QTYPE...] [net SOURCE...]
}
```

- **ZONES** zones it should be authoritative for. If empty, the zones from the configuration block are used.
- **ACTION** (*allow*, *block*, or *filter*) defines the way to deal with DNS queries matched by this rule. The default action is *allow*, which means a DNS query not matched by any rules will be allowed to recurse. The difference between *block* and *filter* is that block returns status code of *REFUSED* while filter returns an empty set *NOERROR*
- **QTYPE** is the query type to match for the requests to be allowed or blocked. Common resource record types are supported. `*` stands for all record types. The default behavior for an omitted `type QTYPE...` is to match all kinds of DNS queries (same as `type *`).
- **SOURCE** is the source IP address to match for the requests to be allowed or blocked. Typical CIDR notation and single IP address are supported. `*` stands for all possible source IP addresses.

## Examples

To demonstrate the usage of plugin acl, here we provide some typical examples.

Block all DNS queries with record type A from 192.168.0.0/16：

~~~ corefile
. {
    acl {
        block type A net 192.168.0.0/16
    }
}
~~~

Filter all DNS queries with record type A from 192.168.0.0/16：

~~~ corefile
. {
    acl {
        filter type A net 192.168.0.0/16
    }
}
~~~

Block all DNS queries from 192.168.0.0/16 except for 192.168.1.0/24:

~~~ corefile
. {
    acl {
        allow net 192.168.1.0/24
        block net 192.168.0.0/16
    }
}
~~~

Allow only DNS queries from 192.168.0.0/24 and 192.168.1.0/24:

~~~ corefile
. {
    acl {
        allow net 192.168.0.0/24 192.168.1.0/24
        block
    }
}
~~~

Block all DNS queries from 192.168.1.0/24 towards a.example.org:

~~~ corefile
example.org {
    acl a.example.org {
        block net 192.168.1.0/24
    }
}
~~~

## Metrics

If monitoring is enabled (via the _prometheus_ plugin) then the following metrics are exported:

- `coredns_acl_blocked_requests_total{server, zone, view}` - counter of DNS requests being blocked.

- `coredns_acl_filtered_requests_total{server, zone, view}` - counter of DNS requests being filtered.

- `coredns_acl_allowed_requests_total{server, view}` - counter of DNS requests being allowed.

The `server` and `zone` labels are explained in the _metrics_ plugin documentation.
