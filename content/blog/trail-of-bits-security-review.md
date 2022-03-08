+++
title = "Trail Of Bits Security Review"
description = "A security review of CoreDNS conducted by Trail of Bits"
tags = ["Security Review", "Threat Model", "Trail of Bits"]
date = "2022-02-24T00:00:00-00:00"
author = "coredns"
+++

Trail of Bits (https://trailofbits.com) conducted a security review and threat model of CoreDNS.

Quoting from the security review summary:

> "The audit uncovered one high-severity issue (TOB-CDNS-8) concerning a bug that could lead to cache poisoning attacks.
> The majority of the other issues are of informational or low severity; these include several resulting from insufficient
> data validation, specifically from assumptions about the data processed by various functions, which we discovered by
> running fuzzing harnesses. Most of the findings pertain to denial-of-service vulnerabilities."

Security Review: https://github.com/trailofbits/publications/blob/master/reviews/coredns-securityreview.pdf

Threat Model: https://github.com/trailofbits/publications/blob/master/reviews/coredns-threatmodel.pdf

At this time, the following PRs have been opened to address issues raised in the report: 

https://github.com/coredns/coredns/pull/5085 (TOB-CDNS-1)
https://github.com/coredns/coredns/pull/5108 (TOB-CDNS-5)
https://github.com/coredns/coredns/pull/5168 (TOB-CDNS-2)
https://github.com/coredns/coredns/pull/5169 (TOB-CDNS-3)
https://github.com/coredns/coredns/pull/5170 (TOB-CDNS-4)
https://github.com/coredns/coredns/pull/5171 (TOB-CDNS-15)
https://github.com/coredns/coredns/pull/5172 (TOB-CDNS-11)
https://github.com/coredns/coredns/pull/5173 (TOB-CDNS-9)
https://github.com/coredns/coredns/pull/5174 (TOB-CDNS-8)
https://github.com/coredns/coredns/pull/5220 (TOB-CDNS-10)
https://github.com/coredns/coredns/pull/5224 (TOB-CDNS-14)
https://github.com/coredns/coredns/pull/5225 (TOB-CDNS-7)
https://github.com/coredns/coredns/pull/5226 (TOB-CDNS-6)
