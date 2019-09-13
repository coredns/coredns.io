+++
title = "ens"
description = "*ens* - serve DNS records from the Ethereum Name Service."
weight = 11
tags = [  "plugin" , "ens" ]
categories = [ "plugin", "external" ]
date = "2019-09-13T09:00:00+01:00"
repo = "https://github.com/wealdtech/coredns-ens"
home = "https://github.com/wealdtech/coredns-ens/blob/master/README.md"
+++

## Description

The *ens* plugin serves DNS records from the Ethereum Name Service.  Ethereum provides an authoritative source of DNS records for relevant domains, allowing authoritative data to be served by any nameserver without it having write-access to the DNS records themselves.

It is recommended that this comes after *rewrite* in the `plugins.cfg` file.

## Syntax

```
ens {
  # connection is the connection to an Ethereum node.  It is *highly*
  # recommended that a local node is used, as remote connections can
  # cause DNS requests to time out.
  # This can be either a path to an IPC socket or a URL to a JSON-RPC
  # endpoint.
  connection CONNECTION

  # ethlinknameservers are the names of the nameservers that serve
  # EthLink domains.  This will usually be the name of this server,
  # plus potentially one or more others.
  ethlinknameservers SERVER...

  # ipfsgatewaya is the address of an ENS-enabled IPFS gateway.
  # This value is returned when a request for an A record of an Ethlink
  # domain is received and the domain has a contenthash record in ENS but
  # no A record.  Multiple values can be supplied, separated by a space,
  # in which case all records will be returned.
  ipfsgatewaya ADDRESS...

  # ipfsgatewayaaaa is the address of an ENS-enabled IPFS gateway.
  # This value is returned when a request for an AAAA record of an Ethlink
  # domain is received and the domain has a contenthash record in ENS but
  # no A record.  Multiple values can be supplied, separated by a space,
  # in which case all records will be returned.
  ipfsgatewayaaaa ADDRESS...
}
```

## Examples

```
ens {
  connection /home/ethereum/.ethereum/geth.ipc
  ethlinknameservers ns1.ethdns.xyz ns2.ethdns.xyz
  ipfsgatewaya 176.9.154.81
  ipfsgatewayaaaa 2a01:4f8:160:4069::2
}
```
