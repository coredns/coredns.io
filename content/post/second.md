+++
Categories = []
Description = ""
Keywords = []
Tags = []
date = "2016-08-02T08:56:42-07:00"
title = "Sane and Simple"
img = "html-code.jpg"
weight = 3
+++
Simple to configure. Sane defaults.

The following `Corefile` snippet shows the configuration for `miek.nl`:

~~~ txt
miek.nl:53 {
    file /var/lib/coredns/miek.nl.signed {
        transfer to *
        transfer to 185.49.141.42
    }
    prometheus
    errors stdout
    log stdout
}
~~~
