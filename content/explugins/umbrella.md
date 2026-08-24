+++
title = "umbrella"
description = "*umbrella* adds Cisco Umbrella identity data to forwarded DNS queries."
weight = 10
tags = [ "plugin", "umbrella" ]
categories = [ "plugin", "external" ]
date = "2026-08-24T18:19:19+00:00"
repo = "https://github.com/xdkr/coredns-umbrella"
home = "https://github.com/xdkr/coredns-umbrella#readme"
+++

## Description

The *umbrella* plugin adds Cisco Umbrella EDNS0 identity data to DNS
queries, enabling internal IP reporting.

It is recommended that this comes before *forward* in the `plugin.cfg` file.

## Syntax

``` txt
umbrella organization_id ORGANIZATION_ID device_id DEVICE_ID
```

* **ORGANIZATION_ID** is an unsigned 32-bit decimal integer.
* **DEVICE_ID** is a 16-character hexadecimal identifier.

## Example

``` corefile
. {
    umbrella organization_id 12345678 device_id 0123456789abcdef
    forward . 208.67.222.222 208.67.220.220
}
```
