+++
title = "k8s_gateway"
description = "*k8s_gateway* - plugin to resolve all types of external Kubernetes resources."
weight = 10
tags = [  "plugin" , "k8s" ]
categories = [ "plugin", "external", "kubernetes" ]
date = "2020-09-19T12:00:00-08:00"
repo = "https://github.com/ori-edge/k8s_gateway"
home = "https://github.com/ori-edge/k8s_gateway/blob/master/README.md"
+++

## Description

This plugin is very similar to [k8s_external](https://coredns.io/plugins/k8s_external/) but supporting all types of Kubernetes external resources - Ingress, Service of type LoadBalancer and `networking.x-k8s.io/Gateway` (when it becomes available). 

This plugin relies on it's own connection to the k8s API server and doesn't share any code with the existing [kubernetes](https://coredns.io/plugins/kubernetes/) plugin. The assumption is that this plugin can now be deployed as a separate instance (alongside the internal kube-dns) and act as a single external DNS interface into your Kubernetes cluster(s).

## Syntax

```
k8s_gateway [ZONE...] 
```

Optionally, you can specify what kind of resources to watch and the default TTL to return in response, e.g.

```
k8s_gateway example.com {
    resources Ingress
    ttl 10
}
```

## Example


``` corefile
. {
  k8s_gateway example.com
}
```

With the above configuration the plugin will behave in the following way:

1. All DNS queries will first be matched against the configured zone - `example.com`
2. If there’s a hit, the next step is to match it against any of the existing Ingress resources. The lookup is performed against FQDNs configured in `spec.rules[*].host` fields of the Ingress. At this stage, the result can be returned to the user with IPs collected from the `.status.loadBalancer.ingress`.
3. If no matching Ingress was found, the search continues with the Services objects. Since services don’t really have domain names, the lookup is performed using the `serviceName.namespace` as the key.
4. If there’s a match, it is returned to the end-user in a similar way, alternatively the plugin responds with NXDOMAIN.

## Supported features

`k8s_gateway` resolves Kubernetes resources with their external IP addresses based on zones specified in the configuration. This plugin will resolve the following type of resources:

| Kind | Matching Against | External IPs are from | 
| ---- | ---------------- | -------- |
| Ingress | all FQDNs from `spec.rules[*].host` matching configured zones | `.status.loadBalancer.ingress` |
| Service[*] | `name.namespace` + any of the configured zones | `.status.loadBalancer.ingress` | 

[*]: Only resolves service of type LoadBalancer

Currently only supports A-type queries, all other queries result in NODATA responses.

>  This plugin is **NOT** supposed to be used for intra-cluster DNS resolution and by default will not contain the default upstream [kubernetes](https://coredns.io/plugins/kubernetes/) plugin.