+++
title = "pocketbase"
description = "*pocketbase* - PocketBase backend for CoreDNS"
weight = 10
tags = [  "plugin" , "pocketbase" ]
categories = [ "plugin", "external" ]
date = "2025-04-11T03:31:26+00:00"
repo = "https://github.com/tinkernels/coredns-pocketbase"
home = "https://github.com/tinkernels/coredns-pocketbase/blob/master/README.md"
+++

# pocketbase

PocketBase backend for CoreDNS

## Name

pocketbase - PocketBase backend for CoreDNS

## Description

This plugin uses PocketBase as a backend to store DNS records. These will then can served by CoreDNS. The backend uses a
simple single table data structure that can add and remove records from the DNS server.

## Syntax

```
pocketbase {
    [listen LISTEN]
    [data_dir DATA_DIR]
    [su_email SU_EMAIL]
    [su_password SU_PASSWORD]
    [default_ttl DEFAULT_TTL]
    [cache_capacity CACHE_CAPACITY]
}
```

- `listen` pocketbase listening http address, default to `[::]:8090`,
- `data_dir` directory to store pocketbase data, default to `pb_data`,
- `su_email` superuser login email, can be overwritten by environment variable `COREDNS_PB_SUPERUSER_EMAIL`, default to `su@pocketbase.internal`,
- `su_password` superuser password, can be overwritten by environment variable `COREDNS_PB_SUPERUSER_PWD`, default to `pwd@pocketbase.internal`,
- `default_ttl` default ttl to use, default to `30`,
- `cache_capacity` zone data cache capacity, `0` to disable cache, default to `0`.

## Features

### Supported Record Types

- A
- AAAA
- CNAME
- SOA
- TXT
- NS
- MX
- CAA
- SRV

*P.S.wildcard records supported*

### Cache

Use `github.com/dgraph-io/ristretto` as in-memory cache handler, handle cache refreshing with PocketBase event subscription mechanism.

## Concept

### PocketBase

[PocketBase](https://github.com/pocketbase/pocketbase) use sqlite3 as storage, and comes with a web console.

This plugin with init a super user and dns model in PocketBase, the admin console with look like

![PocketBase admin console](assets/image/pocketbase-admin.png)

#### Model in PocketBase

```go
type Record struct {
	Zone       string `db:"zone" json:"zone"`               // The DNS zone this record belongs to
	Name       string `db:"name" json:"name"`               // The name of the record (without the zone)
	RecordType string `db:"record_type" json:"record_type"` // The type of DNS record (A, AAAA, TXT, etc.)
	Ttl        uint32 `db:"ttl" json:"ttl"`                 // Time to live for the record in seconds
	Content    string `db:"content" json:"content"`         // The content of the record in JSON format
}
```

### DNS records

DNS records content stored as JSON.

```go
// ARecord represents an A (IPv4) DNS record
type ARecord struct {
	Ip net.IP `json:"ip"` // IPv4 address
}
```
```go
// AAAARecord represents an AAAA (IPv6) DNS record
type AAAARecord struct {
	Ip net.IP `json:"ip"` // IPv6 address
}
```
```go
// TXTRecord represents a TXT DNS record
type TXTRecord struct {
	Text string `json:"text"` // Text content of the record
}
```
```go
// CNAMERecord represents a CNAME DNS record
type CNAMERecord struct {
	Host string `json:"host"` // Target hostname
	Zone string `json:"zone"` // Zone of the record
}
```
```go
// NSRecord represents an NS (Name Server) DNS record
type NSRecord struct {
	Host string `json:"host"` // Name server hostname
}
```
```go
// MXRecord represents an MX (Mail Exchange) DNS record
type MXRecord struct {
	Host       string `json:"host"`       // Mail server hostname
	Preference uint16 `json:"preference"` // Priority of the mail server
}
```
```go
// SRVRecord represents an SRV (Service) DNS record
type SRVRecord struct {
	Priority uint16 `json:"priority"` // Priority of the service
	Weight   uint16 `json:"weight"`   // Weight for load balancing
	Port     uint16 `json:"port"`     // Port number of the service
	Target   string `json:"target"`   // Target hostname
}
```
```go
// SOARecord represents an SOA (Start of Authority) DNS record
type SOARecord struct {
	Ns      string `json:"ns"`      // Primary name server
	MBox    string `json:"mbox"`    // Email address of the administrator
	Refresh uint32 `json:"refresh"` // Refresh interval in seconds
	Retry   uint32 `json:"retry"`   // Retry interval in seconds
	Expire  uint32 `json:"expire"`  // Expiration time in seconds
	MinTtl  uint32 `json:"minttl"`  // Minimum TTL in seconds
}
```
```go
// CAARecord represents a CAA (Certification Authority Authorization) DNS record
type CAARecord struct {
	Flag  uint8  `json:"flag"`  // Critical flag
	Tag   string `json:"tag"`   // Property identifier
	Value string `json:"value"` // Property value
}
```

## Setup (as an external plugin)

Add this as an external plugin in `plugin.cfg` file from CoreDNS repo

```
pocketbase:github.com/tinkernels/coredns-pocketbase
```

*P.S.place pocketbase above cache plugin is recommended.*

Then run

```shell script
$ go generate
$ go build
```

Add any required modules to CoreDNS code as prompted.

## Credits

Inspired by

- [https://github.com/wenerme/coredns-pdsql](https://github.com/wenerme/coredns-pdsql)
- [https://github.com/arvancloud/redis](https://github.com/arvancloud/redis)
- [https://github.com/cloud66-oss/coredns_mysql](https://github.com/cloud66-oss/coredns_mysql)
