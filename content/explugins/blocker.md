+++
title = "blocker"
description = "*blocker* - Domain blocker plugin for CoreDNS."
weight = 10
tags = [  "plugin" , "dns" ]
categories = [ "plugin", "external" ]
date = "2025-04-13T11:58:24+09:00"
repo = "https://github.com/icyflame/blocker"
home = "https://github.com/icyflame/blocker/blob/master/README.org"
+++

## Description

The blocker plugin blocks a list of domains provided in a file written using the AdBlock Plus syntax
format. The list of blocked domains is loaded into memory at start-up, and periodically after that.

## Syntax

~~~ txt
blocker path-to-blocklist-file refresh-interval file-syntax empty-response-type
~~~

1. `path-to-blocklist-file`: Absolute path to a file that contains the list of blocked domains
2. `refresh-interval`: Interval after which the file is read from disk and loaded into memory
   periodically. Duration must be a string which can be parsed by Go's
   [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration).
3. `file-syntax`: One of either `hosts` or `abp`. `hosts` files have a list of blocked domains,
   which will be blocked by exact match. `abp` syntax supports prefixes and multiple subdmoains
   using a single line in the file
4. `empty-response-type`: One of either `empty` or `nxdomain`. `empty` will return `0.0.0.0` or `::`
   as the response when a DNS query for a blocked domain is made. `nxdomain` will return a
   non-existent domain name response for blocked domains.

## Metadata

This plugin exports the metadata key `blocker/request-blocked`. The value of this key will be `YES`
when a domain is blocked and `NO` in all other cases.

If the `metadata` plugin is enabled, then this key can be added to the log line for each query using
the `log` plugin:

``` corefile
.:53 {
  metadata

  log . "{common} {/blocker/request-blocked}"

  blocker /home/user/blocklist_file 1h abp empty

  forward . 1.1.1.1
}
```

## Examples

In this configuration, we block domains that are listed in the `/home/user/blocklist.abp` file and
send an empty response (`0.0.0.0`) when domains in this list are quiered. This file will be reloaded
from disk into memory every hour after start-up:

``` corefile
. {
  blocker /home/user/blocklist.abp 1h abp empty

  forward . 1.1.1.1
}
```
