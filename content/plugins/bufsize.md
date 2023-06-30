+++
title = "bufsize"
description = "*bufsize* sizes EDNS0 buffer size to prevent IP fragmentation."
weight = 7
tags = ["plugin", "bufsize"]
categories = ["plugin"]
date = "2022-05-10T17:23:57.8775785"
+++

## Description
*bufsize* limits a requester's UDP payload size.
It prevents IP fragmentation, mitigating certain DNS vulnerabilities.
This will only affect queries that have an OPT RR (EDNS(0)).

## Syntax
```txt
bufsize [SIZE]
```

**[SIZE]** is an int value for setting the buffer size.
The default value is 1232, and the value must be within 512 - 4096.
Only one argument is acceptable, and it covers both IPv4 and IPv6.

## Examples
Enable limiting the buffer size of outgoing query to the resolver (172.31.0.10):
```corefile
. {
    bufsize 1500
    forward . 172.31.0.10
    log
}
```

Enable limiting the buffer size as an authoritative nameserver:
```corefile
. {
    bufsize 1220
    file db.example.org
    log
}
```

## Considerations
- Setting 1232 bytes to bufsize may avoid fragmentation on the majority of networks in use today, but it depends on the MTU of the physical network links.
