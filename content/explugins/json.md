+++
title: "JSON"
description: "*json* - query a JSON-formatted DNS server"
weight: 10
tags = [ "plugin", "json", "external" ]
categories = [ "external", "plugin" ]
date = "2025-03-05T15:51:45-08:00"
repo = "https://github.com/xinbenlv/coredns-json"
home = "https://github.com/xinbenlv/coredns-json/blob/master/README.md"
+++

## Name

*json* - query a JSON-formatted DNS server

## Description

The *json* plugin queries a JSON-formatted DNS server and returns the result as a DNS response.

## Syntax

```
json <URI>
```

* **URI** (required): The URI of the JSON-formatted DNS server.

## Example

```
. {
    json https://your-json-dns-server.com/api/v1/dns
}
```

## Supported record types

The *json* plugin supports the following DNS record types:

- A
- AAAA
- CNAME
- MX
- TXT

