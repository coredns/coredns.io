+++
title = "flatten"
description = "*<plugin name>* provides minimal CNAME flattening mirroring the Cloudflare RFC 1034 compliant implementation."
weight = 10
tags = [  "plugin" , "<plugin name>" ]
categories = [ "plugin", "external" ]
date = "2017-07-22T12:37:19+01:00"
repo = "https://github.com/litobro/flatten"
home = "https://github.com/litobro/flatten/blob/main/README.md"
+++

## Description

The *<plugin name>* plugin is an attempt to provide a minimal CNAME flattening solution that is RFC compliant in accordance with the Cloudflare implementation.

## Syntax

~~~ txt
flatten [FROM] [TO] [DNSIP:PORT]
~~~

- `FROM`: Original requested NAME
- `TO`: Name to overwrite the A and AAAA records from
- `DNSIP:PORT`: The DNS server and port to resolve the `TO` records from

## Example Corefile
```
example.org:53 {
    log
    flatten example.org google.ca 1.1.1.1:53

    forward . 1.1.1.1
}
```

### Example output
Run the server
```
$ ./coredns -conf Corefile

example.org.:53
CoreDNS-1.11.2
linux/amd64, go1.21.8, 8de4531d-dirty
[INFO] plugin/flatten: 127.0.0.1:37237 - [example.org.] flattened to [google.ca.] via 1.1.1.1:53
[INFO] plugin/flatten: 127.0.0.1:36670 - [example.org.] flattened to [google.ca.] via 1.1.1.1:53
```

Make a DNS request to the server
```
$ nslookup example.org 127.0.0.1

Server:         127.0.0.1
Address:        127.0.0.1#53

Name:   example.org
Address: 142.250.217.67
Name:   example.org
Address: 2607:f8b0:400a:80b::2003
Name:   example.org
Address: 142.250.217.67
Name:   example.org
Address: 2607:f8b0:400a:804::2003

$ dig TXT example.org @127.0.0.1

; <<>> DiG 9.18.18-0ubuntu0.22.04.1-Ubuntu <<>> TXT example.org @127.00.0.1
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 63388
;; flags: qr rd ra ad; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
;; QUESTION SECTION:
;example.org.                   IN      TXT

;; ANSWER SECTION:
example.org.            86400   IN      TXT     "v=spf1 -all"
example.org.            86400   IN      TXT     "6r4wtj10lt2hw0zhyhk7cgzzffhjp7fl"

;; Query time: 99 msec
;; SERVER: 127.0.0.1#53(127.00.0.1) (UDP)
;; WHEN: Wed Apr 17 09:42:46 MDT 2024
;; MSG SIZE  rcvd: 131