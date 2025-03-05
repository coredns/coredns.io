+++
title: "JSON"
description: "The json plugin enables JSON-formatted DNS response output."
weight: 10
+++

# json

## Name

*json* - transforms DNS responses into JSON format

## Description

The *json* plugin converts DNS responses into structured JSON format, making them easier to process by JSON-based tools and pipelines. This is particularly useful for:

- Logging systems that consume JSON
- API responses
- JSON-based analytics pipelines

## Syntax

```
json <URI>
```

* **URI** (required): The URI of the JSON-formatted DNS server.

## Examples

**Basic configuration:**
```corefile
example.org {
    json http://api.your-http-json-dns-server.com/api/v1/
    forward . 8.8.8.8
}
```

The query will be sent to the JSON-formatted DNS server at the given URI with parameters `name` and `type`.

Such as `http://api.your-http-json-dns-server.com/api/v1/?name=example.org&type=A`.


**Sample JSON output:**
```json
{
  "dns": {
    "opcode": "QUERY",
    "rcode": "NOERROR",
    "questions": [
      {
        "name": "example.org.",
        "type": "A"
      }
    ],
    "answers": [
      {
        "name": "example.org.",
        "type": "A",
        "ttl": 3600,
        "data": "93.184.216.34"
      }
    ],
    "timestamp": "2024-02-20T14:30:45Z"
  }
}
```

## Compatibility

* CoreDNS 1.11.0 and later
* Supports the following DNS record types: A, AAAA, CNAME, MX, TXT for now
