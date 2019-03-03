+++
title = "hosts"
description = "*hosts* enables serving zone data from a `/etc/hosts` style file."
weight = 16
tags = [ "plugin", "hosts" ]
categories = [ "plugin" ]
date = "2019-03-03T09:28:16.707194"
+++

## Description

The hosts plugin is useful for serving zones from a `/etc/hosts` file. It serves from a preloaded
file that exists on disk. It checks the file for changes and updates the zones accordingly. This
plugin only supports A, AAAA, and PTR records. The hosts plugin can be used with readily
available hosts files that block access to advertising servers.

The plugin reloads the content of the hosts file every 5 seconds. Upon reload, CoreDNS will use the new definitions.
Should the file be deleted, any inlined content will continue to be served. When the file is restored, it will then again be used.

This plugin can only be used once per Server Block.

## The hosts file

Commonly the entries are of the from `IP_address canonical_hostname [aliases...]` as explained by the hosts(5) man page.

Examples:

~~~
# The following lines are desirable for IPv4 capable hosts
127.0.0.1       localhost
192.168.1.10    example.com            example

# The following lines are desirable for IPv6 capable hosts
::1                     localhost ip6-localhost ip6-loopback
fdfc:a744:27b5:3b0e::1  example.com example
~~~

### PTR records

PTR records for reverse lookups are generated automatically by CoreDNS (based on the hosts file entries) and cannot be created manually.

## Syntax

~~~
hosts [FILE [ZONES...]] {
    [INLINE]
    ttl SECONDS
    no_reverse
    reload DURATION
    fallthrough [ZONES...]
}
~~~

* **FILE** the hosts file to read and parse. If the path is relative the path from the *root*
  directive will be prepended to it. Defaults to /etc/hosts if omitted. We scan the file for changes
  every 5 seconds.
* **ZONES** zones it should be authoritative for. If empty, the zones from the configuration block
   are used.
* **INLINE** the hosts file contents inlined in Corefile. If there are any lines before fallthrough
   then all of them will be treated as the additional content for hosts file. The specified hosts
   file path will still be read but entries will be overridden.
* `ttl` change the DNS TTL of the records generated (forward and reverse). The default is 3600 seconds (1 hour).
* `reload` change the period between each hostsfile reload. A time of zero seconds disable the feature. Examples of valid durations: "300ms", "1.5h" or "2h45m" are valid duration with units "ns" (nanosecond), "us" (or "µs" for microsecond), "ms" (millisecond), "s" (second), "m" (minute), "h" (hour).
* `no_reverse` disable the automatic generation of the `in-addr.arpa` or `ip6.arpa` entries for the hosts
* `fallthrough` If zone matches and no record can be generated, pass request to the next plugin.
  If **[ZONES...]** is omitted, then fallthrough happens for all zones for which the plugin
  is authoritative. If specific zones are listed (for example `in-addr.arpa` and `ip6.arpa`), then only
  queries for those zones will be subject to fallthrough.

## Examples

Load `/etc/hosts` file.

~~~ corefile
. {
    hosts
}
~~~

Load `example.hosts` file in the current directory.

~~~
. {
    hosts example.hosts
}
~~~

Load example.hosts file and only serve example.org and example.net from it and fall through to the
next plugin if query doesn't match.

~~~
. {
    hosts example.hosts example.org example.net {
        fallthrough
    }
}
~~~

Load hosts file inlined in Corefile.

~~~
. {
    hosts example.hosts example.org {
        10.0.0.1 example.org
        fallthrough
    }
}
~~~

## See also

The form of the entries in the `/etc/hosts` file are based on IETF [RFC 952](https://tools.ietf.org/html/rfc952) which was updated by IETF [RFC 1123](https://tools.ietf.org/html/rfc1123).
