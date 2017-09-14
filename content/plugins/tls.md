+++
title = "tls"
description = "*tls* allows you to configure the server certificates for the TLS and gRPC servers. For other types of servers it is ignored."
weight = 26
tags = [ "plugin", "tls" ]
categories = [ "plugin" ]
date = "2017-09-10T18:11:52.766707"
+++

## Syntax

~~~ txt
tls CERT KEY CA
~~~

## Examples

Start a DNS-over-TLS server.

~~~
tls://.:4453 {
	tls cert.pem key.pem ca.pem
	proxy . /etc/resolv.conf
}
~~~

Start a DNS-over-gRPC server. If the `tls` directive were omitted, then
it would use plain HTTP not HTTPS.

~~~
grpc://.:443 {
	tls cert.pem key.pem ca.pem
	proxy . /etc/resolv.conf
}
~~~
