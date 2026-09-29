+++
title = "Getting Started with RFC 2136 Dynamic DNS"
description = "Configure authenticated, persistent DNS updates for a small local zone."
tags = ["Documentation", "Tutorial", "DDNS"]
date = "2026-09-29T00:00:00Z"
author = "houyuwushang"
+++

The *dynupdate* plugin accepts RFC 2136 DNS UPDATE requests for a writable
authoritative zone. A DHCP updater can use it to maintain a device's forward
and reverse DNS records without rewriting a zone file or restarting CoreDNS.
CoreDNS provides the DNS service, not DHCP address allocation or lease management.

> Availability: as of September 29, 2026, *dynupdate* is merged into CoreDNS
> `master` but is not included in v1.14.7. Use a build from `master` for this
> tutorial. The initial implementation is experimental and intended for small,
> infrequently updated unsigned zones.

This walkthrough uses loopback and port 1053, so it needs neither a container
nor permission to bind port 53. The shell commands assume a Unix-like shell
and BIND's `tsig-keygen`, `nsupdate`, and `dig` tools.

## Build CoreDNS

Use the Go version required by the repository's `go.mod`:

~~~ sh
git clone https://github.com/coredns/coredns.git
cd coredns
go build -o coredns .
./coredns -plugins
~~~

The plugin list must include `dynupdate` and `tsig`. An older release
binary will not recognize the `dynupdate` directive.

## Create a Private Working Directory and Key

From the source directory, create a working directory for the configuration,
seed zones, key, and databases:

~~~ sh
umask 077
mkdir ddns-demo
cd ddns-demo
tsig-keygen -a hmac-sha256 dhcp-key.example.test. > update.key
~~~

The key file is shared by CoreDNS and the update client. Keep it private;
do not put it in source control or reuse a public example secret. TSIG
authenticates messages but does not encrypt their contents.

## Seed the Forward and Reverse Zones

Create `example.test.zone` with:

~~~ zone
$ORIGIN example.test.
@ 60 IN SOA ns.example.test. hostmaster.example.test. 1 3600 600 86400 60
@ 60 IN NS ns.example.test.
ns 60 IN A 192.0.2.53
~~~

Create `2.0.192.in-addr.arpa.zone` with:

~~~ zone
$ORIGIN 2.0.192.in-addr.arpa.
@ 60 IN SOA ns.example.test. hostmaster.example.test. 1 3600 600 86400 60
@ 60 IN NS ns.example.test.
~~~

These are initialization files, not the live database. CoreDNS never writes
updates back to them. Once a database is initialized, that database is
authoritative; editing the seed will not replace acknowledged updates.

## Configure the Server

Create `Corefile` in the same directory:

~~~ corefile
example.test:1053 {
    bind 127.0.0.1
    errors
    tsig {
        secrets update.key
        require_opcode UPDATE
    }
    dynupdate {
        file example.test.zone
        database example.test.db
        allow dhcp-key.example.test. host.example.test. A AAAA DHCID ANY
    }
}

2.0.192.in-addr.arpa:1053 {
    bind 127.0.0.1
    errors
    tsig {
        secrets update.key
        require_opcode UPDATE
    }
    dynupdate {
        file 2.0.192.in-addr.arpa.zone
        database 2.0.192.in-addr.arpa.db
        allow dhcp-key.example.test. 10.2.0.192.in-addr.arpa. PTR DHCID ANY
    }
}
~~~

Each writable zone has its own server block and database. `require_opcode UPDATE`
rejects unsigned updates while leaving ordinary queries unsigned. The `allow`
rules authorize only the named records, not every name in either zone.

The extra `DHCID` and `ANY` grants are useful for DHCP updaters such as Kea D2.
`ANY` permits deleting **all RRsets at that owner name**, not just the other
listed types. Do not put unrelated static records at these authorized names.
For a client that only adds and deletes individual A or PTR records, omit the
permissions it does not need.

Start CoreDNS in the foreground:

~~~ sh
../coredns -conf Corefile
~~~

Keep it running, and run the remaining commands from a second terminal in the
same `ddns-demo` directory. The directory must be writable by the CoreDNS
process. A missing database is initialized on the first query, transfer, or
authenticated update after startup; keep the seed files until then.

## Add and Query Records

Submit signed updates over UDP:

~~~ sh
nsupdate -k update.key <<'EOF'
server 127.0.0.1 1053
zone example.test.
prereq nxrrset host.example.test. A
update add host.example.test. 60 A 192.0.2.10
send
zone 2.0.192.in-addr.arpa.
prereq nxrrset 10.2.0.192.in-addr.arpa. PTR
update add 10.2.0.192.in-addr.arpa. 60 PTR host.example.test.
send
EOF
~~~

Each `send` is a separate transaction. Forward and reverse updates are not
atomic together: a rejected reverse update does not roll back the forward
update. The prerequisites ensure this example does not replace an existing
RRset. Running it again without deleting the records will fail those checks.

Query CoreDNS directly:

~~~ sh
dig @127.0.0.1 -p 1053 host.example.test. A +short
dig @127.0.0.1 -p 1053 -x 192.0.2.10 +short
~~~

The answers should be `192.0.2.10` and `host.example.test.` respectively.
The in-tree *cache* plugin bypasses dynamic zones when configured, but external
recursive resolvers may retain older answers until their TTL expires.

## Verify Persistence and Removal

Stop CoreDNS with Ctrl-C in the first terminal, then start it again with the
same command, directory, and configuration. Repeat both `dig` commands: the
records should still be present even though the seed files do not contain them.
Without `database`, updates would be lost on restart or Corefile reload.

Remove the two records using TCP, selected by `nsupdate -v`:

~~~ sh
nsupdate -v -k update.key <<'EOF'
server 127.0.0.1 1053
zone example.test.
prereq yxrrset host.example.test. A 192.0.2.10
update delete host.example.test. A
send
zone 2.0.192.in-addr.arpa.
prereq yxrrset 10.2.0.192.in-addr.arpa. PTR host.example.test.
update delete 10.2.0.192.in-addr.arpa. PTR
send
EOF
~~~

Both `dig` commands should now produce no answer. To inspect the response code
as well, replace `+short` with `+noall +comments +answer`.

For an authentication check, run the addition command without `-k update.key`.
The server should return `REFUSED`, and the records should remain absent.

## Connect a DHCP Updater

A DHCP deployment replaces the manual `nsupdate` commands with its DNS update
client. Configure that client with the authoritative server address, port,
forward and reverse zones, and the same TSIG key name, algorithm, and secret.
This loopback example accepts only clients on the same machine. For a remote
DHCP updater, bind to the intended private interface and restrict network
access in addition to using TSIG; do not expose the service publicly.

Match authorization to the client's actual messages. Kea D2, for example, can
write DHCID records in both zones and delete all records at a name. If you
expand this example to more lease addresses, add the corresponding owner
permissions. A name of `*` grants access to the entire zone and should only be
used for an updater trusted to manage all of it.

Keep the DHCP client's ownership-conflict checks enabled. A TSIG key identifies
the updater, not the device that owns a lease. The DHCP service remains
responsible for lease expiry, record cleanup, and retrying or reconciling
partial forward/reverse failures. It must actively send deletions; CoreDNS
does not derive DHCP lease expiry from DNS TTLs.

The CoreDNS integration suite exercises a real Kea D2 process with IPv4 and
IPv6 forward/reverse updates, renewal, DHCID conflicts, removal, and name reuse.
Its lease-change notifications are synthetic, so this is not a test of DHCP
address allocation or automatic lease expiry in your network.

## Operational Boundaries

Use a local filesystem with working locks and sync semantics. The database
belongs to one zone and one CoreDNS process; do not share it between independent
servers. Stop CoreDNS before copying a database for an offline backup or
restoring it. Never edit, replace, or delete a live database.

Updates rebuild a bounded zone snapshot and, with persistence, wait for a disk
commit before success is returned. This is not a high-throughput DHCP backend.
IXFR history, automatic DNSSEC signing, multi-primary replication, and
transactions spanning zones are outside the initial implementation.

See the [dynupdate reference][dynupdate] for resource limits, recovery, AXFR and
NOTIFY behavior, and the full permission syntax. The [tsig reference][tsig]
describes key configuration. Implementation and review history are linked from
[the DDNS issue][issue].

[dynupdate]: https://github.com/coredns/coredns/blob/master/plugin/dynupdate/README.md
[tsig]: https://coredns.io/plugins/tsig/
[issue]: https://github.com/coredns/coredns/issues/6254
