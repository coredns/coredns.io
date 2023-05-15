+++
title = "meshname"
description = "*meshname* allows to resolve decentralized .meshname domains"
weight = 10
tags = [  "plugin" , "meshname" , "meship" ]
categories = [ "plugin", "external" ]
date = "2023-05-15T12:00:00+00:00"
repo = "https://github.com/zhoreeq/coredns-meshname"
home = "https://github.com/zhoreeq/meshname"
+++

## Background

Currently, a strict hierarchy is followed when resolving DNS names. It is centrally designed. To address the problem of centrality, there are several approaches. One of these approaches is Meshname. With meshname, the IPv6 address of the authoritative name server responsible for the meshname domain is already encoded in the domain name. Thus, when a meshname domain is to be resolved, the first thing that is done is to decode the encoded authoritative DNS server from the domain. Then the DNS request is sent to this DNS server. Thus, the resolution happens without the typical querying of the hierarchy of DNS servers. Only a connection to the encoded DNS server must exist in order to resolve a meshname domain.
Meshname domains fulfill the criteria of "Decentralized" and "Secure" but not of "Freely Selectable" of Zooko's triangle. However, the "Secure" aspect requires that the connection to the DNS server can be made securely.
Changing the IP address of the authoritative DNS server requires changing the domain name. Furthermore, only one authoritative DNS server can be specified in a meshname domain. A meshname domain cannot be resolved if this server is offline.
A specification of the protocol can be found at [https://github.com/zhoreeq/meshname/blob/master/protocol.md](https://github.com/zhoreeq/meshname/blob/master/protocol.md).

## Description

The *meshname* plugin allows to resolve decentralized .meshname domains. These are domains in which the IPv6 of the authoritative server is decoded. The advantage is that no central instance is needed to resolve the authoritative server. The disadvantage is that the meshname names can look quite ugly.

## Syntax

```
meshname
```

## Example

```
meshname. {
  meshname
}
```

