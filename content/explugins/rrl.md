+++
 title = "rrl"
 description = "*rrl* - mitigate amplification attacks with response rate limiting"
 weight = 10
 tags = [  "plugin" , "rrl" ]
 categories = [ "plugin", "external" ]
 date = "2021-10-07T00:00:00+00:00"
 repo = "https://github.com/coredns/rrl"
 home = "https://github.com/coredns/rrl/blob/master/README.md"
 +++

 ## Description

The *rrl* plugin tracks response rates per category of response. 
The category of a given response consists of the following:

* Prefix of the client IP (per the ipv4/6-prefix-length)
* Requested name (qname) excluding response type of error (see response type below)
* Requested type (qtype) excluding response type of error (see response type below)
* Response type (each corresponding to the configurable per-second allowances)
  * response - for positive responses that contain answers
  * nodata - for NODATA responses
  * nxdomain - for NXDOMAIN responses
  * referrals - for referrals or delegations
  * error - for all DNS errors (except NXDOMAIN)


To better protect against attacks using invalid requests, requested name
and type are not categorized separately for error type requests. In other
words, all error responses are limited collectively per client, regardless
of qname or qtype.

Each category has an account balance which is credited at a rate of the
configured *per-second* allowance for that response type, and debited each
time a response in that category would be sent to a client.  When an account
balance is negative, responses in the category are dropped until the balance
goes non-negative.  Account balances cannot be more positive than  *per-second*
allowance, and cannot be more negative than *window* * *per-second* allowance.

The response rate limiting implementation intends to replicate the behavior
of BIND 9's response rate limiting feature.

When limiting requests, the category of each request is determined by the
prefix of the client IP (per the ipv4/6-prefix-length).

 ## Syntax

~~~txt
rrl [ZONES...] {
    window SECONDS
    ipv4-prefix-length LENGTH
    ipv6-prefix-length LENGTH
    responses-per-second ALLOWANCE
    nodata-per-second ALLOWANCE
    nxdomains-per-second ALLOWANCE
    referrals-per-second ALLOWANCE
    errors-per-second ALLOWANCE
    requests-per-second ALLOWANCE
    max-table-size SIZE
    report-only
}
~~~

* `window SECONDS` - the rolling window in **SECONDS** during which response rates are tracked. Default 15.

* `ipv4-prefix-length LENGTH` - the prefix **LENGTH** in bits to use for identifying a ipv4 client. Default 24.

* `ipv6-prefix-length LENGTH` - the prefix **LENGTH** in bits to use for identifying a ipv6 client. Default 56.

* `responses-per-second ALLOWANCE` - the number of positive responses allowed per second. An **ALLOWANCE** of 0 disables rate limiting of positive responses. Default 0.

* `nodata-per-second ALLOWANCE` - the number of `NODATA` responses allowed per second. An **ALLOWANCE** of 0 disables rate limiting of NODATA responses. Defaults to responses-per-second.

* `nxdomains-per-second ALLOWANCE` - the number of `NXDOMAIN` responses allowed per second. An **ALLOWANCE** of 0 disables rate limiting of NXDOMAIN responses. Defaults to responses-per-second.

* `referrals-per-second ALLOWANCE` - the number of referral responses allowed per second. An **ALLOWANCE** of 0 disables rate limiting of referral responses. Defaults to responses-per-second.

* `errors-per-second ALLOWANCE` - the number of error responses allowed per second (excluding NXDOMAIN). An **ALLOWANCE** of 0 disables rate limiting of error responses. Defaults to responses-per-second.

* `requests-per-second ALLOWANCE` - the number of requests allowed per second. An **ALLOWANCE** of 0 disables rate limiting of requests. Default 0.

* `max-table-size SIZE` - the maximum number of responses to be tracked at one time. When exceeded, rrl stops rate limiting new responses. Defaults to 100000.

* `report-only` -  Do not drop requests/responses when rates are exceeded, only log metrics. Defaults to false.

## Metrics

If monitoring is enabled (via the *prometheus* plugin) then the following metric are exported:

* `coredns_rrl_responses_exceeded_total{client_ip}` - Counter of responses exceeding QPS limit.
* `coredns_rrl_requests_exceeded_total{client_ip}` - Counter of requests exceeding QPS limit.

## Examples

Example 1

~~~ corefile

. {
  rrl . {
    responses-per-second 10
  }
}

~~~

## Bugs / Known Issues / Limitations

BIND9's implementation of Response Rate Limiting will rate limit all wildcard generated records in one account per the base domain of the wild card.  e.g. Both `a.dom.com.` and  `b.dom.com.` would be accounted for as `dom.com.`, if they are generated from the wildcard record `*.dom.com.`

Per the BIND 9.11 ARM...

> Responses generated from local wildcards are counted and limited as if they were for the parent domain name. 
> This controls flooding using random.wild.example.com.

In CoreDNS *rrl* wildcard responses are accounted for individually.

## See Also

[A Quick Introduction to Response Rate Limiting](https://kb.isc.org/docs/aa-01000)

