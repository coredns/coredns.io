+++
title = "kubeforward"
description = "*kubeforward* - dynamically updates DNS forwarders based on changes to a specified Kubernetes Service."
weight = 10
tags = ["plugin", "kubeforward", "kubernetes", "dns"]
categories = ["plugin", "external"]
date = "2025-04-04T17:15:00+00:00"
repo = "https://github.com/deckhouse/coredns-kubeforward"
home = "https://github.com/deckhouse/coredns-kubeforward/blob/main/README.md"
+++

## Description

The **kubeforward** plugin enables CoreDNS to dynamically update its list of DNS forwarders by monitoring changes to a specified Kubernetes Service. It observes `EndpointSlices` associated with the Service and adjusts the DNS forwarding configuration in real-time as endpoints are added, removed, or updated. This functionality enhances the reliability and resilience of DNS services within Kubernetes clusters.

## Syntax

```txt
kubeforward {
    namespace <namespace>
    service_name <service_name>
    port_name <port_name>
    expire <duration>
    health_check <duration>
    force_tcp
    prefer_udp
}
```

## Configuration Parameters

- `namespace` (required): Specifies the Kubernetes namespace where the target Service resides.

- `service_name` (required): The name of the Service to monitor for endpoint changes.

- `port_name`: The name of the port in the Service resource responsible for handling DNS queries.

- `expire`: Duration after which cached connections expire. Default is 10s.

- `health_check`: Interval for performing health checks on upstream servers. Default is 0.5s.

- `force_tcp`: Forces the use of TCP for forwarding queries.

- `prefer_udp`: Prefers the use of UDP for forwarding queries.

## Examples

```
.:53 {
    errors
    log
    kubeforward {
        namespace kube-system
        service_name kube-dns
        port_name dns
        expire 10m
        health_check 5s
        prefer_udp
        force_tcp
    }
}
```

## Limitations

Limited Support for Forward Plugin Options: The plugin utilizes the functionality of the forward plugin for serving DNS under the hood but does not support the full list of classic forward options due to the lack of a public interface for configuring options.



