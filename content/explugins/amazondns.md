+++
title = "amazondns"
description = "*amazondns* - enables serving an authoritative name server using Amazon DNS Server as the backend."
weight = 10
tags = [ "plugin" , "amazondns" ]
categories = [ "plugin", "external" ]
date = "2018-03-06T20:32:00+09:00"
repo = "https://github.com/wadahiro/coredns-amazondns"
home = "https://github.com/wadahiro/coredns-amazondns/blob/master/README.md"
+++

## Description

The *amazondns* plugin behaves **Authoritative name server** using [Amazon DNS Server](https://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_DHCP_Options.html#AmazonDNS) as the backend.

The Amazon DNS server is used to resolve the DNS domain names that you specify in a private hosted zone in Route 53. However, the server acts as **Caching name server**. Although CoreDNS has [proxy plugin](https://github.com/coredns/coredns/tree/master/plugin/proxy) and we can configure Amazon DNS server as the backend, it can't be Authoritative name server. In my case, Authoritative name server is required to handle delegated responsibility for the subdomain. That's why I created this plugin. 

## Syntax

~~~ txt
amazondns ZONE [ADDRESS] {
    soa RR
    ns RR
    nsa RR
}
~~~

* **ZONE** the zone scope for this plugin.
* **ADDRESS** defines the Amazon DNS server address specifically.
  If no **ADDRESS** entry, this plugin resovles it automatically using [Instance Metadata](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html).
* **soa** **RR** SOA record with [RFC 1035](https://tools.ietf.org/html/rfc1035#section-5) style.
* **ns** **RR** NS record(s) with [RFC 1035](https://tools.ietf.org/html/rfc1035#section-5) style.
* **nsa** **RR** A record(s) for the NS(s) with [RFC 1035](https://tools.ietf.org/html/rfc1035#section-5) style.
  The IP address will be the private IP address of the EC2 instance on which CoreDNS is running with this plugin.
  Note: You need to boot CoreDNS on an EC2 instance in the VPC because we can't access to Amazon DNS server from outside the VPC.

## Examples

Create an authoritative name server for `sub.example.org` with two name servers(`ns1.sub.example.org` and `ns2.sub.example.org`).

~~~ txt
. {
    amazondns sub.example.org {
        soa "sub.example.org 3600 IN SOA ns1.sub.example.org hostmaster.sub.example.org (2018030619 3600 1200 1209600 900)"
        ns "sub.example.org 3600 IN NS ns1.sub.example.org"
        ns "sub.example.org 3600 IN NS ns2.sub.example.org"
        nsa "ns1.sub.example.org 3600 IN A 192.168.0.10"
        nsa "ns2.sub.example.org 3600 IN A 192.168.0.130"
    }
}
~~~

