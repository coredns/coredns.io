+++
title = "meship"
description = "*meship* allows to resolve decentralized .meship domains"
weight = 10
tags = [  "plugin" , "meshname" , "meship" ]
categories = [ "plugin", "external" ]
date = "2023-05-15T12:00:00+00:00"
repo = "https://github.com/zhoreeq/coredns-meship"
home = "https://github.com/zhoreeq/meshname"
+++

## Background

Currently, a strict hierarchy is followed when resolving DNS names. It is centrally designed. To address the problem of centrality, there are several approaches. One of these approaches is Meshname.
With a meship domain, the address to which the domain is to be resolved is encoded in the domain name. Thus, when a meship domain is to be resolved, the domain name is decoded first and then returned as a AAAA record.
Meship domains meet the criteria of "Decentralized" and "Secure" but not of "Freely Selectable" of Zooko's triangle.
However, meship domains are not scalable because only one and the same IP address can be resolved at a time. Furthermore, subdomains are not possible. Furthermore, if the IP address is changed, the entire domain name must be changed.
A specification of the protocol can be found at https://github.com/zhoreeq/meshname/blob/master/protocol.md.

## Description

The *meship* plugin allows to resolve decentralized .meship domains. The AAAA record is already decoded in the domain name. This is then returned accordingly.

## Syntax

```
meship
```

## Example

```
meship. {
  meship
}
```

