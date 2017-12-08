+++
title = "ipin"
description = "*ipin* returns IP address and port based on you domain name."
weight = 10
tags = [  "plugin" , "ipin" ]
categories = [ "plugin", "external" ]
date = "2017-12-09T02:00:00+08:00"
repo = "https://github.com/wenerme/wps"
home = "https://github.com/wenerme/wps/blob/master/coredns/plugin/ipin/README.md"
+++

# ipin

*ipin* returns IP address and port based on you domain name. Your IP address is returned
 in the additional section as either an A or AAAA record.

The reply always has an empty answer section. The port are included in the additional
section as a SRV record.

~~~ txt
._port.qname. 0 IN SRV 0 0 <port> .
~~~



## Syntax

~~~ txt
ipin
~~~

## Examples

Start a server on the default port and load the *ipin* plugin.

~~~ corefile
. {
    ipin
}
~~~

When queried for "192-168-1-1.example.org A", CoreDNS will respond with:

~~~ txt
;; QUESTION SECTION:
;192-168-1-1.example.org.	IN	A

;; ADDITIONAL SECTION:
192-168-1-1.example.org. 0	IN	A	192.168.1.1
~~~

When queried for "127-0-0-1-8080.example.org A", CoreDNS will respond with:

~~~ txt
;; QUESTION SECTION:
;127-0-0-1-8080.example.org.	IN	A

;; ADDITIONAL SECTION:
127-0-0-1-8080.example.org. 0	IN	A	127.0.0.1
_port.127-0-0-1-8080.example.org. 0 IN	SRV	0 0 8080 .
~~~
