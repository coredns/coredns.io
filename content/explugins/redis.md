+++
title = "redis"
description = "*redis* - enables reading zone data from redis database."
weight = 10
tags = [  "plugin" , "redis" ]
categories = [ "plugin", "external" ]
date = "2017-10-18T18:01:20+03:30"
repo = "https://github.com/hawell/redis"
home = "https://github.com/hawell/redis/blob/master/README.md"
+++

## Description

*redis* enables reading zone data from redis database.
this plugin should be located right next to *etcd* in *plugins.cfg*

## Syntax

~~~
redis
~~~

redis loads authoritative zones from redis server


address will default to local redis serrver (localhost:6379)
~~~
redis {
    address ADDR
    password PWD
    prefix PREFIX
    suffix SUFFIX
    connect_timeout TIMEOUT
    read_timeout TIMEOUT
    ttl TTL
}
~~~




* `address` is redis server address to connect in the form of *host:port* or *ip:port*.
* `password` is redis server *auth* key
* `connect_timeout` time in ms to wait for redis server to connect
* `read_timeout` time in ms to wait for redis server to respond
* `ttl` default ttl for dns records, 300 if not provided
* `prefix` add PREFIX to all redis keys
* `suffix` add SUFFIX to all redis keys

## Examples

~~~ corefile
. {
    redis example.com {
        address localhost:6379
        password foobared
        connect_timeout 100
        read_timeout 100
        ttl 360
        prefix _dns:
    }
}
~~~

## Reverse Zones

reverse zones is not supported yet.

## Proxy

proxy is not supported yet.

## Zone Format in redis db

### Zones

each zone is stored in redis as a hash map with *zone* as key

~~~
redis-cli>KEYS *
1) "example.com."
2) "example.net."
redis-cli>
~~~

### DNS RRs

dns RRs are stored in redis as json strings inside a hash map using address as field key.
*@* is used for zone's own RR values.

#### A

~~~json
{
    "a":{
        "ip4" : "1.2.3.4",
        "ttl" : 360
    }
}
~~~

#### AAAA

~~~json
{
    "aaaa":{
        "ip6" : "::1",
        "ttl" : 360
    }
}
~~~

#### CNAME

~~~json
{
    "cname":{
        "host" : "x.example.com.",
        "ttl" : 360
    }
}
~~~

#### TXT

~~~json
{
    "txt":{
        "text" : "this is a text",
        "ttl" : 360
    }
}
~~~

#### NS

~~~json
{
    "ns":{
        "host" : "ns1.example.com.",
        "ttl" : 360
    }
}
~~~

#### MX

~~~json
{
    "mx":{
        "host" : "mx1.example.com",
        "priority" : 10,
        "ttl" : 360
    }
}
~~~

#### SRV

~~~json
{
    "srv":{
        "host" : "sip.example.com.",
        "port" : 555,
        "priority" : 10,
        "weight" : 100,
        "ttl" : 360
    }
}
~~~

#### SOA

~~~json
{
    "soa":{
        "ttl" : 100,
        "mbox" : "hostmaster.example.com.",
        "ns" : "ns1.example.com.",
        "refresh" : 44,
        "retry" : 55,
        "expire" : 66
    }
}
~~~

#### Example

~~~
$ORIGIN example.net.
 example.net.                 300 IN  SOA   <SOA RDATA>
 example.net.                 300     NS    ns1.example.net.
 example.net.                 300     NS    ns2.example.net.
 *.example.net.               300     TXT   "this is a wildcard"
 *.example.net.               300     MX    10 host1.example.net.
 sub.*.example.net.           300     TXT   "this is not a wildcard"
 host1.example.net.           300     A     5.5.5.5
 _ssh.tcp.host1.example.net.  300     SRV   <SRV RDATA>
 _ssh.tcp.host2.example.net.  300     SRV   <SRV RDATA>
 subdel.example.net.          300     NS    ns1.subdel.example.net.
 subdel.example.net.          300     NS    ns2.subdel.example.net.
~~~

above zone data should be stored at redis as follow:

~~~
redis-cli> hgetall example.net.
 1) "_ssh._tcp.host1"
 2) "{\"srv\":[{\"ttl\":300, \"target\":\"tcp.example.com.\",\"port\":123,\"priority\":10,\"weight\":100}]}"
 3) "*"
 4) "{\"txt\":[{\"ttl\":300, \"text\":\"this is a wildcard\"}],\"mx\":[{\"ttl\":300, \"host\":\"host1.example.net.\",\"preference\": 10}]}"
 5) "host1"
 6) "{\"a\":[{\"ttl\":300, \"ip\":\"5.5.5.5\"}]}"
 7) "sub.*"
 8) "{\"txt\":[{\"ttl\":300, \"text\":\"this is not a wildcard\"}]}"
 9) "_ssh._tcp.host2"
10) "{\"srv\":[{\"ttl\":300, \"target\":\"tcp.example.com.\",\"port\":123,\"priority\":10,\"weight\":100}]}"
11) "subdel"
12) "{\"ns\":[{\"ttl\":300, \"host\":\"ns1.subdel.example.net.\"},{\"ttl\":300, \"host\":\"ns2.subdel.example.net.\"}]}"
13) "@"
14) "{\"soa\":{\"ttl\":300, \"minttl\":100, \"mbox\":\"hostmaster.example.net.\",\"ns\":\"ns1.example.net.\",\"refresh\":44,\"retry\":55,\"expire\":66},\"ns\":[{\"ttl\":300, \"host\":\"ns1.example.net.\"},{\"ttl\":300, \"host\":\"ns2.example.net.\"}]}"
redis-cli>
~~~
