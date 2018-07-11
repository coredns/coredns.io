+++
title = "dnssec"
description = "*dnssec* enable on-the-fly DNSSEC signing of served data."
weight = 7
tags = [ "plugin", "dnssec" ]
categories = [ "plugin" ]
date = "2018-07-11T10:14:28.428601"
+++

## Description

With *dnssec* any reply that doesn't (or can't) do DNSSEC will get signed on the fly. Authenticated
denial of existence is implemented with NSEC black lies. Using ECDSA as an algorithm is preferred as
this leads to smaller signatures (compared to RSA). NSEC3 is *not* supported.

This plugin can only be used once per Server Block.

## Syntax

~~~
dnssec [ZONES... ] {
    key file KEY...
    cache_capacity CAPACITY
}
~~~

The specified key is used for all signing operations. The DNSSEC signing will treat this key as a
CSK (common signing key), forgoing the ZSK/KSK split. All signing operations are done online.
Authenticated denial of existence is implemented with NSEC black lies. Using ECDSA as an algorithm
is preferred as this leads to smaller signatures (compared to RSA). NSEC3 is *not* supported.

If multiple *dnssec* plugins are specified in the same zone, the last one specified will be
used (See [bugs](#bugs)).

* **ZONES** zones that should be signed. If empty, the zones from the configuration block
    are used.

* `key file` indicates that **KEY** file(s) should be read from disk. When multiple keys are specified, RRsets
  will be signed with all keys. Generating a key can be done with `dnssec-keygen`: `dnssec-keygen -a
  ECDSAP256SHA256 <zonename>`. A key created for zone *A* can be safely used for zone *B*. The name of the
  key file can be specified in one of the following formats

    * basename of the generated key `Kexample.org+013+45330`
    * generated public key `Kexample.org+013+45330.key`
    * generated private key `Kexample.org+013+45330.private`

* `cache_capacity` indicates the capacity of the cache. The dnssec plugin uses a cache to store
  RRSIGs. The default for **CAPACITY** is 10000.

## Metrics

If monitoring is enabled (via the *prometheus* directive) then the following metrics are exported:

* `coredns_dnssec_cache_size{server, type}` - total elements in the cache, type is "signature".
* `coredns_dnssec_cache_hits_total{server}` - Counter of cache hits.
* `coredns_dnssec_cache_misses_total{server}` - Counter of cache misses.

The label `server` indicated the server handling the request, see the *metrics* plugin for details.

## Examples

Sign responses for `example.org` with the key "Kexample.org.+013+45330.key".

~~~ corefile
example.org {
    dnssec {
        key file Kexample.org.+013+45330
    }
    whoami
}
~~~

Sign responses for a kubernetes zone with the key "Kcluster.local+013+45129.key".

~~~
cluster.local {
    kubernetes
    dnssec {
      key file Kcluster.local+013+45129
    }
}
~~~
