+++
title = "clouddns"
description = "*clouddns* enables serving zone data from GCP Cloud DNS."
weight = 11
tags = [ "plugin", "clouddns" ]
categories = [ "plugin" ]
date = "2019-12-12T16:13:55.315947"
+++

## Description

The *clouddns* plugin is useful for serving zones from resource record
sets in GCP Cloud DNS. This plugin supports all [Google Cloud DNS
records](https://cloud.google.com/dns/docs/overview#supported_dns_record_types). This plugin can
be used when CoreDNS is deployed on GCP or elsewhere. Note that this plugin accesses the resource
records through the Google Cloud API. For records in a privately hosted zone, it is not necessary to
place CoreDNS and this plugin in the associated VPC network. In fact the private hosted zone could
be created without any associated VPC and this plugin could still access the resource records under
the hosted zone.

## Syntax

~~~ txt
clouddns [ZONE:PROJECT_ID:HOSTED_ZONE_NAME...] {
    credentials [FILENAME]
    fallthrough [ZONES...]
}
~~~

*   **ZONE** the name of the domain to be accessed. When there are multiple zones with overlapping
    domains (private vs. public hosted zone), CoreDNS does the lookup in the given order here.
    Therefore, for a non-existing resource record, SOA response will be from the rightmost zone.

*   **PROJECT_ID** the project ID of the Google Cloud project.

*   **HOSTED_ZONE_NAME** the name of the hosted zone that contains the resource record sets to be
    accessed.

*   `credentials` is used for reading the credential file.

*   **FILENAME** GCP credentials file path (normally a .json file).

*   `fallthrough` If zone matches and no record can be generated, pass request to the next plugin.
    If **[ZONES...]** is omitted, then fallthrough happens for all zones for which the plugin is
    authoritative. If specific zones are listed (for example `in-addr.arpa` and `ip6.arpa`), then
    only queries for those zones will be subject to fallthrough.

*   **ZONES** zones it should be authoritative for. If empty, the zones from the configuration block

## Examples

Enable clouddns with implicit GCP credentials and resolve CNAMEs via 10.0.0.1:

~~~ txt
example.org {
    clouddns example.org.:gcp-example-project:example-zone
    forward . 10.0.0.1
}
~~~

Enable clouddns with fallthrough:

~~~ txt
example.org {
    clouddns example.org.:gcp-example-project:example-zone example.com.:gcp-example-project:example-zone-2 {
        fallthrough example.gov.
    }
}
~~~

Enable clouddns with multiple hosted zones with the same domain:

~~~ txt
. {
    clouddns example.org.:gcp-example-project:example-zone example.com.:gcp-example-project:other-example-zone
}
~~~
