+++
title = "Cure53 Security Assessment"
description = "Details of the CoreDNS security assessment by Cure53."
tags = ["Security"]
date = "2018-03-15T09:34:29+00:00"
author = "coredns"
enabled = "default"
+++

Being an [incubating CNCF](https://www.cncf.io/projects/) project makes us eligible for nice things
like a security assessment (cue ominous music).

The CNCF asked [Cure53](https://cure53.de) to perform such an assessment.

TL;DR: CoreDNS is in good shape, but Cure53 did find one critical issue (which we've fixed with the
CoreDNS 1.1.1 release):

> ### DNS-01-003 Cache: DNS Cache poisoning via malicious Response (Critical)
>
> The CoreDNS application allows to configure the caching of the DNS responses via the
> cache plugin. It was discovered that CoreDNS only verifies the transaction IDs but fails
> to check whether the domain in a request matches the response. This can be abused to
> inject malicious A records in the cache of the DNS server.
> As the CoreDNS application has a different cache for each domain

The other three issues found will be tracked via github issues, like
[plugin/rewrite: log bypass](https://github.com/coredns/coredns/issues/1610), and
[plugin/secondary: Denial-of-Service via endless Zone Transfer](plugin/secondary: Denial-of-Service
via endless Zone Transfer). Third one was a generic DDoS.

On a positive note the final report includes quotes like these:

> The CoreDNS software tested by Cure53 during this March 2018 assessment has made
a clearly positive impression.

<!-- -->

> To conclude, even though four issues were found during this Cure53 assessment, they
were generally - with a single exception - minor, miscellaneous and manageable.
Despite Cure53 testers’ considerable efforts, the software was found to be hard to
corrupt. Therefore, the CoreDNS project stands out as secure, robust and legitimately
security-aware.

The full report can be found [here](/assets/DNS-01-report.pdf). As for future improvements in
CoreDNS: we will increase the use of fuzzing, increase test coverage  and look closer at DNS DoS
mitigations, such as DNS Cookies (described in [RFC 7873](https://tools.ietf.org/html/rfc7873)).
