+++
title = "responsefilter"
description = "*responsefilter* filters DNS responses based on FQDN and IP CIDR blocklists to protect against DNS spoofing."
weight = 10
tags = ["plugin", "responsefilter"]
categories = ["plugin", "external"]
date = "2026-01-22T17:00:00+01:00"
repo = "https://github.com/isovalent/responsefilter"
home = "https://github.com/isovalent/responsefilter"
+++

## Description

The *responsefilter* plugin inspects DNS responses from upstream servers and blocks responses where the returned IP address matches a configured blocklist for specific domains. When a blocked response is detected, CoreDNS returns a REFUSED status instead of the spoofed IP address.

This plugin helps protect against DNS spoofing attacks and malicious DNS responses by allowing administrators to define which IP ranges are not acceptable for specific domains.

## Syntax

```
responsefilter {
    block DOMAIN CIDR [CIDR...]
}
```

* **DOMAIN** - the domain name to apply the filter to (supports subdomains)
* **CIDR** - one or more IP CIDR ranges to block for this domain

**Important:** The responsefilter directive must be placed before the forward directive in your Corefile.

## Examples

Block specific IP ranges for a domain:

```
.:53 {
    responsefilter {
        block abc.com 10.1.1.0/24
    }
    forward . 8.8.8.8
}
```

Block multiple CIDR ranges for multiple domains:

```
.:53 {
    responsefilter {
        block abc.com 10.1.1.0/24 192.168.0.0/16
        block xyz.com 172.16.0.0/12
    }
    forward . 8.8.8.8
}
```
