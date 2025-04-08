+++
title = "recursor"
description = "recursor resolves domains using defined IP addresses or resolving other mapped domains using defined resolvers"
weight = 10
tags = [  "plugin" , "recursor" ]
categories = [ "plugin", "external" ]
date = "2023-01-30T16:15:00+02:00"
repo = "https://github.com/kinjelom/coredns-recursor"
home = "https://github.com/kinjelom/coredns-recursor/blob/main/README.md"
+++

## Description

The recursor resolves domains using defined IP addresses or resolving other mapped domains using defined resolvers.

![](https://github.com/kinjelom/coredns-recursor/raw/main/docs/flow.png)

## Syntax

~~~ txt
recursor {   
    [external-yaml config-file-path]
    [external-json config-file-path]

    [verbose 0..2]
    [resolver non-default {
        urls udp://ip-1:port udp://ip-n:port
        [timeout_ms 500]
    }]
    [alias alias-name | * {
        [hosts domain-1 domain-n]
        [ips ip-1 ip-n]
        [resolver_name non-default]
        [ttl custom-ttl]
    }]
}
~~~

The `recursor` definition:
- `verbose` - stdout logs level:
  - `0` - minimal
  - `1` - middle
  - `2` - talkative
- `resolvers` - other DNS servers:
  - *map-key/id*: name of resolver, `default` overrides system default resolver
  - `urls`: list of URL addresses, example: `udp://127.0.0.1:53` (system default is represented by `://default`)
  - `timeout_ms`: resolver connection timeout in millisecods
- `aliases` - domain aliases:
  - *map-key/id*: name of alias, subdomain or `*` if you want the recursor to be a DNS repeater
  - `ips`: IP addresses to return as a part of answer
  - `hosts`: domains to be resolved so that the IP addresses obtained in this way will be returned as a part of answer
  - `resolver_name`: the defined resolver reference, default is... `default` of course :)
  - `ttl`: DNS record Time To Live in seconds


## Metrics

If monitoring is enabled (via the *prometheus* directive) [the following metric is exported](https://github.com/kinjelom/coredns-recursor/blob/main/metrics.go).
- [Definition](https://github.com/kinjelom/coredns-recursor/blob/main/metrics.go)
- [Grafana Dashboard](https://github.com/kinjelom/coredns-recursor/blob/main/docs/dashboard.json)

![](https://github.com/kinjelom/coredns-recursor/blob/main/docs/dashboard.png?raw=true)

## Health

This plugin implements dynamic health checking. It will always return healthy though.

## Examples

#### Corefile

```txt
recursor {   
    [external-yaml config-file-path]
    [external-json config-file-path]

    [verbose 0..2]
    resolver dns-c {
        urls udp://1.1.1.1:53 udp://1.0.0.1:53
        timeout_ms 500
    }
    resolver dns-g {
        urls udp://8.8.8.8:53 udp://8.8.4.4:53
    }  
    resolver demo {
        urls udp://10.0.0.1:53
    }  
    alias alias1 {
        hosts www.example.org www.example.com
        resolver_name dns-c
        ttl 11
    }
    alias alias2 {
        ips 10.0.0.1 10.0.0.2
        ttl 12
    }
    alias alias3 {
        ips 10.0.0.1 10.0.0.2
        hosts www.example.net
        resolver_name dns-g
        ttl 13
    }
    alias alias4 {
        hosts www.example.net
        ttl 14
    }  

    alias * {
        resolver_name demo
        ttl 15
    }  
}
```

#### External YAML

```yaml
resolvers:
  dns-c:
    urls: [ udp://1.1.1.1:53, udp://1.0.0.1:53 ]
    timeout_ms: 500
  dns-g:
    urls: [ udp://8.8.8.8:53, udp://8.8.4.4:53 ]
  demo:
    urls: [ udp://10.0.0.1:53 ]
aliases:
  alias1:
    hosts: [ www.example.org, www.example.com ]
    resolver_name: dns-c
    ttl: 11
  alias2:
    ips: [ 10.0.0.1, 10.0.0.2 ]
    ttl: 12
  alias3:
    ips: [ 10.0.0.1, 10.0.0.2 ]
    hosts: [ www.example.net ]
    resolver_name: dns-g
    ttl: 13
  alias4:
    hosts: [ www.example.net ]
    ttl: 14
  "*":
    resolver_name: demo
    ttl: 15
```

#### External JSON

```json
{
  "resolvers": {
    "dns-c": {
      "urls": [ "udp://1.1.1.1:53", "udp://1.0.0.1:53" ],
      "timeout_ms": 500
    },
    "dns-g": {
      "urls": [ "udp://8.8.8.8:53", "udp://8.8.4.4:53" ]
    },
    "demo": {
      "urls": [ "udp://10.0.0.1:53" ]
    }
  },
  "aliases": {
    "alias1": {
      "hosts": [ "www.example.org", "www.example.com" ],
      "resolver_name": "dns-c",
      "ttl": 11
    },
    "alias2": {
      "ips": [ "10.0.0.1", "10.0.0.2" ],
      "ttl": 12
    },
    "alias3": {
      "ips": [ "10.0.0.1", "10.0.0.2" ],
      "hosts": [ "www.example.net" ],
      "resolver_name": "dns-g",
      "ttl": 13
    },
    "alias4": {
      "hosts": [ "www.example.net" ],
      "ttl": 14
    },
    "*": {
      "resolver_name": "demo",
      "ttl": 15
    }
  }
}
```

## Also See

See the [manual](https://coredns.io/manual).
