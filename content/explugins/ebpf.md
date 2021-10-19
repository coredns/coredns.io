+++
title = "ebpf"
description = "*ebpf* - attach an eBPF XDP program to a specified interface."
weight = 10
tags = [  "plugin" , "ebpf" ]
categories = [ "plugin", "external" ]
date = "2021-10-14T00:00:00+00:00"
repo = "https://github.com/InfobloxOpen/ebpf"
home = "https://github.com/InfobloxOpen/ebpf/blob/master/README.md"
+++

## Description

This *experimental* plugin allows you to use an eBPF XDP program to analyze and filter traffic before it reaches CoreDNS,
and report very basic Prometheus metrics. When CoreDNS exits, the program will be detached from the interface. 

This generic solution serves in part as an example of how you can integrate an eBPF XDP program with CoreDNS with a
custom plugin. But due to the generic nature, map entry is somewhat cryptic and metrics must be defined in the Corefile,
limiting their scope. When writing your own plugin, you can tailor it to work with a specific XDP program, for example,
to enable easier human-readable data entry or publish more advanced metrics.

## Syntax

~~~ txt
ebpf {
  elf PROGRAM
  if INTERFACE
  map [KEY] VALUE
  metric NAME KEY POS LEN "HELP"
}
~~~

* `elf` **PROGRAM** - the ELF program to attach.  See notes below on program requirements.
* `if` **INTERFACE** - the interface to attach to
* `map` **KEY** **VALUE** - the hexidecimal string representations of the **KEY** and **VALUE** of
  an entry to load into the eBPF map. You may specify the `map` option more than once to add multiple
  items to the map. If **KEY** is not specified, the entry is treated as an array value.  To make multi-field
  values easier to visually digest, **VALUE** may be delimited by dots.  e.g. `012345678.0000000000000000.9ABCDEF0`
  This is for legibility of the Corefile only; any dots in **VALUE** are ignored by the parser.  When *debug* is used
  the values written to log are not delimited.
* `metric` **NAME** **KEY** **POS** **LEN** "**HELP**" - when used in conjunction with the *prometheus* plugin, register
  a Prometheus "gauge" metric to expose a eBPF map value as an integer metric. The metric is named **NAME** with help
  text of **HELP**.  The map value to use is determined by the **KEY**, byte position **POS**, and length **LEN** in
  bytes.  **LEN** can be at most 8 bytes (64 bits).  The integer value should be little endian.
  
Please be aware of the considerable footgun potential of this plugin.  An XDP program attached to an interface will act
on _all_ ingress packets to the interface - not just packets bound for CoreDNS.

## eBPF Program and Map Requirements

The program must be an XDP program, and main function named `xdp_prog`.
The map must be named `xdp_map`.

Some example programs written in C are included in https://github.com/InfobloxOpen/ebpf/tree/master/example_programs.

## Examples

If `my_xdp_program.o` defines a map with a 4 byte key, and the following struct as a value ...
```
struct maprec {
  __be32  ip4net;  // ipv4 network
  __be32  ip4mask; // ipv4 mask
  __be32  count;   // packet count
};
```

The following will attach `my_xdp_program.o` to `eth0`, and load data for IP network `10.11.0.0` ,
IP mask `255.255.0.0`, and a count of zero (`0A0B0000`, `FFFF0000`, and `00000000` respectively) into key `00000000` of
the map.

```
. {
  ebpf {
    if eth0
    elf my_xdp_program.o
    map 00000000 0A0B0000FFFF000000000000
  }
}
```
The following adds dots to the map value to make it easier to read.

```
. {
  ebpf {
    if eth0
    elf my_xdp_program.o
    map 00000000 0A0B0000.FFFF0000.00000000
  }
}
```

The following will enable debug to monitor map values and log when they change.

```
. {
  debug
  ebpf {
    if eth0
    elf my_xdp_program.o
    map 00000000 0A0B0000.FFFF0000.00000000
  }
}
```

The following adds map entries without specifying keys.  Each map entry is inserted as an array value, with an 
automatically incrementing key.

```
. {
  ebpf {
    if eth0
    elf my_xdp_program.o
    map 0A0B0000.FFFF0000.00000000
    map 0A0C0000.FFFF0000.00000000
    map 0A0D0000.FFFF0000.00000000
  }
}
```

The example above is equivalent to the following but with keys specified.  Note that the keys are little endian in
this example.

```
. {
  ebpf {
    if eth0
    elf my_xdp_program.o
    map 00000000 0A0B0000.FFFF0000.00000000
    map 01000000 0A0C0000.FFFF0000.00000000
    map 02000000 0A0D0000.FFFF0000.00000000
  }
}
```

The following exposes a Prometheus metric.  The metric is named `coredns_ebpf_example_total` and the value will reflect
the rightmost 4 bytes from map entry `02000000`.

```
. {
  prometheus :9153
  ebpf {
    if eth0
    elf my_xdp_program.o
    map 00000000 0A0B0000.FFFF0000.00000000
    map 01000000 0A0C0000.FFFF0000.00000000
    map 02000000 0A0D0000.FFFF0000.00000000
    metric example_total 02000000 8 4 "Example count."
  }
}
```
