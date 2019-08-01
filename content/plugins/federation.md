+++
title = "federation"
description = "*federation* enables federated queries to be resolved via the kubernetes plugin."
weight = 14
tags = [ "plugin", "federation" ]
categories = [ "plugin" ]
date = "2019-08-01T14:00:49.171397"
+++

## Description

Enabling this plugin allows
[Federated](https://kubernetes.io/docs/tasks/federation/federation-service-discovery/) queries to be
resolved via the kubernetes plugin.

Enabling *federation* without also having *kubernetes* is a noop.

## Syntax

~~~
federation [ZONES...] {
    NAME DOMAIN
}
~~~

* Each **NAME** and **DOMAIN** defines federation membership. One entry for each. A duplicate
  **NAME** will silently overwrite any previous value.

## Examples

Here we handle all service requests in the `prod` and `stage` federations.

~~~
. {
    kubernetes cluster.local
    federation cluster.local {
        prod prod.feddomain.com
        staging staging.feddomain.com
    }
    forward . 192.168.1.12
}
~~~
