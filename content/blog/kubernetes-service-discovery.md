+++
date = "2016-11-08T07:33:08Z"
description = "Using CoreDNS for service discovery in Kubernetes."
tags = ["Kubernetes", "Service", "Discovery", "Kube-DNS", "Documentation"]
title = "CoreDNS for Kubernetes Service Discovery"
author = "miek"
+++

[Infoblox](https://www.infoblox.com/)'s John Belamaric has published a
blog post on
[how to use CoreDNS instead of
kube-dns](https://community.infoblox.com/t5/Community-Blog/CoreDNS-for-Kubernetes-Service-Discovery/ba-p/8187)
in Kubernetes.

A little excerpt:

> Kubernetes includes a DNS server, Kube-DNS, for use in service discovery. This DNS server utilizes
> the libraries from SkyDNS to serve DNS requests for Kubernetes pods and services. The author of
> SkyDNS2, Miek Gieben, has a new DNS server, CoreDNS, that is built with a more modular, extensible
> framework. Infoblox has been working with Miek to adapt this DNS server as an alternative to
> Kube-DNS.

Read more on their
[site](https://community.infoblox.com/t5/Community-Blog/CoreDNS-for-Kubernetes-Service-Discovery/ba-p/8187).
