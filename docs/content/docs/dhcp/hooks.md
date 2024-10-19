---
title: "Hooks"
---

## Hook Methods

### `onDHCPRequestBefore(request)`

Called before the DHCP response is generated.

- `request`: See https://pkg.go.dev/beryju.io/gravity/pkg/roles/dhcp#Request4

### `onDHCPRequestAfter(request, response)`

Called after the DHCP response is generated.

- `request`: See https://pkg.go.dev/beryju.io/gravity/pkg/roles/dhcp#Request4
- `response`: See https://pkg.go.dev/github.com/insomniacslk/dhcp/dhcpv4#DHCPv4

## Environment

{{< gravity-hook-env >}}
