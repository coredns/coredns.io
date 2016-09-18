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

The following `Corefile` snippet shows the configuration for `example.org`:

~~~ txt
example.org:53 {
    file /var/lib/coredns/example.org.signed {
        transfer to *
    }
    prometheus
    errors stdout
    log stdout
}
~~~
