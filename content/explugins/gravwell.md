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
    Ingest-Secret IngestSecretToken
    Cleartext-Target 192.168.1.1:4023
    Tag dns
    Encoding json
    Log-Level INFO
    #Cleartext-Target 192.168.1.2:4023 #second indexer
    #Ciphertext-Target 192.168.1.1:4024
    #Insecure-Novalidate-TLS true #disable TLS certificate validation
    #Ingest-Cache-Path /tmp/coredns_ingest.cache #enable the local ingest cache
    #Max-Cache-Size-MB 1024
}
~~~

* **Ingest-Secret** defines the token used to authenticate with indexers.  **Ingest-Secret** is required.
* **Cleartext-Target** defines the address and port for a remote indexer using a TCP connection.  IPv4 and IPv6 addresses as well as host names are supported.
* **Ciphertext-Target** defines the address and port for a remote indexer using a TLS connection.  IPv4 and IPv6 addresses as well as host names are supported.
* **Tag** specifies the tag that DNS audit logs are assigned.  Can be any alphanumeric value without special characters or spaces.  A valid Tag value is required.
* **Encoding** specifies the format of transported DNS audit logs.  Options are _json_ or _text_.  Deafult is _json_.
* **Insecure-Novalidate-TLS** toggles certificate validation on TLS connections.  Validation is on by default.
* **Log-Level** specifies the logging verbosity over the integrated gravwell tag.  Options are _OFF_ _INFO_ _WARN_ _ERROR_.  Default is _ERROR_.
* **Ingest-Cache-Path** specifies a file path for the cache system which engages when indexer connectivity is lost.  Path must be an absolute path to a writable file.
* **Max-Cache-Size-MB** specifies in megabytes the maximum size of the cache file.  This is used as a safty net.  Zero value is the default and represents unlimited.

## Examples

### No local cache with single indexer over TCP

A sample Corefile which sends DNS requests to a single indexer over an unencrypted connection.  Local cache is disabled.

~~~
gravwell {
    Ingest-Secret IngestSecretToken
    Cleartext-Target 192.168.1.1:4023
    Tag dns
  }
~~~

### TLS connection to two indexers with no TLS validation

A sample Corefile which sends DNS requests to two indexers over a TLS connection and accepts unsigned certificates. Local cache is disabled.
IPv4 and IPv6 addresses are supported for both the Cleartext and Ciphertext targets.  IPv6 addresses must be enclosed in brackets.

~~~
gravwell {
    Ingest-Secret IngestSecretToken
    Ciphertext-Target 192.168.1.1:4024
    Ciphertext-Target [fe80::dead:beef:feed:febe%p1p1]:4024 #connecting to link local address via device p1p1
    Tag dns
    Encoding json
    Log-Level INFO
  }
~~~

### TLS connection to two indexers with no TLS validation

A sample Corefile which sends DNS requests to two indexers over a TLS connection and accepts unsigned certificates. Local cache is disabled.

~~~
gravwell {
    Ingest-Secret IngestSecretToken
    Ciphertext-Target 192.168.1.1:4024
    Ciphertext-Target [dead::beef]:4024
    Insecure-Novalidate-TLS true
    Tag dns
    Encoding json
    Log-Level INFO
  }
~~~

### Local cache for high reliability operation

A sample Corefile which sends DNS requests to two indexers and enables a local cache should indexer communication fail.  Up to 1GB of data can be locally cached.

~~~
gravwell {
    Ingest-Secret IngestSecretToken
    Cleartext-Target 192.168.1.1:4023
    Ciphertext-Target 192.168.1.2:4024
    Insecure-Novalidate-TLS true
    Ingest-Cache-Path /tmp/coredns_ingest.cache
    Max-Cache-Size-MB 1024
    Tag dns
    Encoding json
    Log-Level INFO
  }
~~~

## See Also

[Getting started](https://dev.gravwell.io/docs/#!quickstart/community-edition.md) with Gravwell Community Edition
[Community Edition Licenses](https://www.gravwell.io/activate-community-edition)
[Ingest API and code](https://github.com/gravwell/ingest)
