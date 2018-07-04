+++
title = "gravwell"
description = "*gravwell* - integrate into Gravwell auditing."
weight = 10
tags = [  "plugin" , "gravwell" ]
categories = [ "plugin", "external" ]
date = "2018-07-04T20:25:00+00:00"
repo = "https://github.com/gravwell/coredns"
home = "https://github.com/gravwell/coredns/blob/master/README.md"
+++

## Description

This plugin allows for directly integrating DNS auditing into Gravwell. The plugin acts as an
integrated ingester and ships DNS requests and responses directly to a Gravwell instance.

DNS Requests and responses can be encoded as text, JSON, or as a packed binary format.

## Syntax

~~~
gravwell {
    Ingest-Secret=IngestSecretToken
    Cleartext-Target=192.168.1.1:4023
    Tag=dns
    Encoding=json
    Log-Level=INFO
    #Cleartext-Target=192.168.1.2:4023 #second indexer
    #Ciphertext-Target=192.168.1.1:4024
    #Insecure-Novalidate-TLS=true #disable TLS certificate validation
    #Ingest-Cache-Path=/tmp/coredns_ingest.cache #enable the local ingest cache
    #Max-Cache-Size-MB=1024
  }
~~~

TODO: explain options

* **FROM** is the base domain to match for the request to be resolved. If not specified the zones
  from the server block are used.

## Examples

TODO

## See Also
