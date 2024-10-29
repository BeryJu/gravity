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

### `dhcp` Object

#### `Opt(code: number, data: byte[])`

Create a DHCP option with the code and data given.

## Examples

### Set Option 43 for UniFi Adoption

```javascript
const UniFiPrefix = [0x01, 0x04];
const UniFiIP = net.parseIP("192.168.1.100", "v4");
function onDHCPRequestAfter(req, res) {
    res.UpdateOption(dhcp.Opt(43, [...UniFiPrefix, ...UniFiIP]))
}
```
