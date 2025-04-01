+++
title = "gslb"
description = "*gslb* (Global Server Load Balancing) functionality in CoreDNS."
weight = 10
tags = [  "plugin" , "gslb" ]
categories = [ "plugin", "external" ]
date = "2025-02-09T17:00:00+02:00"
repo = "https://github.com/dmachard/coredns-gslb"
home = "https://github.com/dmachard/coredns-gslb#readme"
+++

## Description

This plugin provides support for GSLB, enabling advanced load balancing and failover mechanisms based on backend health checks and policies. 
It is particularly useful for managing geographically distributed services or for ensuring high availability and resilience.

### Features:
- IPv4 and IPv6 support
- **Health Checks**:
  - HTTPS
  - TCP
  - ICMP
- **Selection Modes**:
  - **Failover**: Routes traffic to the highest-priority available backend.
  - **Random**: Distributes traffic randomly across backends.
  - **Round Robin**: Cycles through backends in sequence.

## Syntax

~~~
gslb DB_YAML_FILE [ZONES...] {
    max_stagger_start "120s"
    resolution_idle_timeout "3600s"
    batch_size_start 100
}
~~~

* **DB_YAML_FILE** The GSLB configuration file in YAML format. If the path is relative, the path from the *root*
  plugin will be prepended to it.
* **ZONES** Specifies the zones the plugin should be authoritative for. If not provided, the zones from the CoreDNS configuration block are used.


## Examples


Load the `gslb.example.com` zone from `db.gslb.example.com` and enable GSLB records on it

~~~ corefile
. {
    file db.gslb.example.com
    gslb gslb_config.example.com.yml db.gslb.example.com {
        max_stagger_start "120s"
        resolution_idle_timeout "3600s"
        batch_size_start 100
    }
}
~~~

Where `db.gslb.example.com` would contain 

~~~ text
$ORIGIN gslb.example.com.
@       3600    IN      SOA     ns1.example.com. admin.example.com. (
                                2024010101 ; Serial
                                7200       ; Refresh
                                3600       ; Retry
                                1209600    ; Expire
                                3600       ; Minimum TTL
                                )
        3600    IN      NS      ns1.gslb.example.com.
        3600    IN      NS      ns2.gslb.example.com.
~~~

And `gslb_config.example.com.yml` would contain 

~~~ yaml
records:
  webapp.gslb.example.com.:
    mode: "failover"
    record_ttl: 30
    scrape_interval: 10s
    backends:
    - address: "172.16.0.10"
      priority: 1
      healthchecks:
      - type: http
        params:
          port: 443
          uri: "/"
          host: "localhost"
          expected_code: 200
          enable_tls: true
    - address: "172.16.0.11"
      priority: 2
      healthchecks:
      - type: http
        params:
          port: 443
          uri: "/"
          host: "localhost"
          expected_code: 200
          enable_tls: true
~~~
