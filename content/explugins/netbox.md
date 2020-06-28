+++
title = "netbox"
description = "*netbox* - enables reading zone data from a netbox instance."
weight = 10
tags = [  "plugin" , "netbox" ]
categories = [ "plugin", "external" ]
date = "2020-06-28T15:01:20+01:00"
repo = "https://github.com/oz123/coredns-netbox-plugin/"
home = "https://github.com/oz123/coredns-netbox-plugin/README.md"
+++

## Description

*netbox* enables reading zone data from a [netbox][1] instance.

## Syntax

~~~
netbox {
  url  http://10.0.0.2:9000/api/ipam/ip-addresses
	token youSekretAPITokenForNetbox
}
~~~

The plugin will delegate search to the next plugin if a record isn't found.
If a record is found a record is sent and the query processing is stopped.

[1]: https://netbox.readthedocs.io/en/stable/
