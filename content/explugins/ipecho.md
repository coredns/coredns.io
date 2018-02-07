+++
title = "ipecho"
description = "*ipecho* - parses the IP out of a subdomain and echos it back as an record"
weight = 180
tags = [  "plugin" , "ipecho" ]
categories = [ "ipecho", "external" ]
date = "2017-09-29T16:18:35+01:00"
repo = "http://github.com/Eun/coredns-ipecho"
home = "https://github.com/Eun/coredns-ipecho/blob/master/README.md"
+++

## Description

*ipecho* parses the IP out of a subdomain and echos it back as an record.

## Example
```
A IN 127.0.0.1.example.com. -> A: 127.0.0.1
AAAA IN ::1.example.com. -> AAAA: ::1
```

## Syntax
```
ipecho {
    domain example1.com
    domain example2.com
    ttl 2629800
}
```

* **domain** adds the domain that should be handled
* **ttl** defines the ttl that should be used in the response
* **debug** enables debug logging
