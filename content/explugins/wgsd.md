+++
title = "wgsd"
description = "*wgsd* - serves WireGuard peer information via DNS-SD semantics to enable endpoint discovery and NAT traversal."
weight = 10
tags = [  "plugin" , "wgsd" ]
categories = [ "plugin", "external" ]
date = "2020-05-27T17:05:00-08:00"
repo = "https://github.com/jwhited/wgsd"
home = "https://github.com/jwhited/wgsd/blob/master/README.md"
+++

## Description

*wgsd* serves [WireGuard](https://www.wireguard.com) peer information via DNS-SD ([RFC6763](https://tools.ietf.org/html/rfc6763)) semantics. This enables dynamic discovery of WireGuard endpoint addressing (both IP address and port number) with the added benefit of NAT-to-NAT WireGuard connectivity where [UDP hole punching](https://en.wikipedia.org/wiki/UDP_hole_punching) is supported.

Following RFC6763 this plugin provides a listing of peers via PTR records at the namespace `_wireguard._udp.<zone>`. The target for the PTR records is `<base32PubKey>._wireguard._udp.<zone>` which corresponds to SRV records. SRV targets are of the format `<base32PubKey>.<zone>`. When querying the SRV record for a peer, the target A/AAAA records will be included in the "additional" section of the response. Public keys are represented in Base32 rather than Base64 to allow for their use in node names where they are treated as case-insensitive by the DNS.

## Syntax

```
wgsd ZONE DEVICE
```

## Example

This configuration:
```
$ cat Corefile
.:5353 {
  wgsd example.com. wg0
}
```

With the following WireGuard peers:
```
$ sudo wg show
interface: wg0
  public key: JeZlz14G8tg1Bqh6apteFCwVhNhpexJ19FDPfuxQtUY=
  private key: (hidden)
  listening port: 51820

peer: xScVkH3fUGUv4RrJFfmcqm8rs3SEHr41km6+yffAHw4=
  endpoint: 203.0.113.1:7777
  allowed ips: 10.0.0.1/32
  latest handshake: 14 hours, 24 minutes, 40 seconds ago
  transfer: 840.64 KiB received, 85.54 KiB sent

peer: syKB97XhGnvC+kynh2KqQJPXoOoOpx/HmpMRTc+r4js=
  endpoint: 198.51.100.1:8888
  allowed ips: 10.0.0.2/32
  latest handshake: 4 days, 15 hours, 8 minutes, 12 seconds ago
  transfer: 1.38 MiB received, 139.42 KiB sent
```

Will respond with:
```
$ dig @127.0.0.1 -p 5353 _wireguard._udp.example.com. PTR +noall +answer +additional
_wireguard._udp.example.com. 0	IN	PTR	yutrled535igkl7bdlerl6m4vjxsxm3uqqpl4nmsn27mt56ad4ha====._wireguard._udp.example.com.
_wireguard._udp.example.com. 0	IN	PTR	wmrid55v4enhxqx2jstyoyvkicj5pihkb2tr7r42smiu3t5l4i5q====._wireguard._udp.example.com.
$
$ dig @127.0.0.1 -p 5353 yutrled535igkl7bdlerl6m4vjxsxm3uqqpl4nmsn27mt56ad4ha====._wireguard._udp.example.com. SRV +noall +answer +additional
yutrled535igkl7bdlerl6m4vjxsxm3uqqpl4nmsn27mt56ad4ha====._wireguard._udp.example.com. 0	IN SRV 0 0 7777 yutrled535igkl7bdlerl6m4vjxsxm3uqqpl4nmsn27mt56ad4ha====.example.com.
yutrled535igkl7bdlerl6m4vjxsxm3uqqpl4nmsn27mt56ad4ha====.example.com. 0	IN A 203.0.113.1
$
$ dig @127.0.0.1 -p 5353 wmrid55v4enhxqx2jstyoyvkicj5pihkb2tr7r42smiu3t5l4i5q====._wireguard._udp.example.com. SRV +noall +answer +additional
wmrid55v4enhxqx2jstyoyvkicj5pihkb2tr7r42smiu3t5l4i5q====._wireguard._udp.example.com. 0	IN SRV 0 0 8888 wmrid55v4enhxqx2jstyoyvkicj5pihkb2tr7r42smiu3t5l4i5q====.example.com.
wmrid55v4enhxqx2jstyoyvkicj5pihkb2tr7r42smiu3t5l4i5q====.example.com. 0	IN A 198.51.100.1
```

Converting public keys to Base64 with coreutils:
```
$ echo yutrled535igkl7bdlerl6m4vjxsxm3uqqpl4nmsn27mt56ad4ha==== | tr '[:lower:]' '[:upper:]' | base32 -d | base64
xScVkH3fUGUv4RrJFfmcqm8rs3SEHr41km6+yffAHw4=
$ echo wmrid55v4enhxqx2jstyoyvkicj5pihkb2tr7r42smiu3t5l4i5q==== | tr '[:lower:]' '[:upper:]' | base32 -d | base64
syKB97XhGnvC+kynh2KqQJPXoOoOpx/HmpMRTc+r4js=
```