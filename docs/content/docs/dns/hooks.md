---
title: "Hooks"
---

## Hook Methods

### `onDNSRequestBefore(request)`

Called before the DNS response is generated.

- `request`: See https://pkg.go.dev/beryju.io/gravity/pkg/roles/dns/utils#DNSRequest

### `onDNSRequestAfter(request, response)`

Called after the DNS response is generated.

- `request`: See https://pkg.go.dev/beryju.io/gravity/pkg/roles/dns/utils#DNSRequest
- `response`: See https://pkg.go.dev/github.com/miekg/dns#Msg

## Environment

{{< gravity-hook-env >}}
