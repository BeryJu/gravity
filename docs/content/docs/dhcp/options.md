---
title: Options
weight: 5
---

## Tag

### `tagName`:

The `tagName` corresponds to a human-readable name for a DHCP option tag; any of these tag names can be used:

- `subnet_mask`
- `router`
- `time_server`
- `name_server`
- `domain_name`
- `bootfile`
- `tftp_server`

*Conflicts with `tag`*

### `tag`:

This field allows for setting the raw DHCP option tag to any value.

*Conflicts with `tagName`*

## Value

### `value`:

Set the value this option should be set to. When used in conjunction with `tagName`, Gravity will automatically encode the value correctly. When used with `tag`, the value needs to be correctly escaped manually.

Example:

```yaml
- tagName: router
  value: 10.1.2.3
```

*Conflicts with `value64` and `valueHex`*

### `value64`:

Set the value this option should be set to, encoded in base64. This allows for pre-encoding the value when using `tag`, and representing data that can't be encoded in an ASCII string. This value should be set to an array of base64-strings, all of which are concatenated.

*Conflicts with `value`*

### `valueHex`:

Set the value this option should be set to, encoded in Hexadecimal. This value should be set to an array of hex-strings, all of which are concatenated.

Example:

```yaml
- tag: 43
  valueHex:
    - 0104C0A8030A
```

*Conflicts with `value`*

## Defaults

Gravity applies some default options when not explicitly configured in the scope settings, but will always prefer user-configured settings if available.

### Subnet Mask

- Tag name: `subnet_mask`
- Tag: `1`

This option defaults to the subnet mask from the CIDR configured for the scope.

### DNS Server

- Tag name: `name_server`
- Tag: `6`

This option defaults to the IP address of the Gravity instance responding to a DHCPREQUEST.

### Hostname

- Tag: `12`

This option defaults to the hostname provided by the client in the DHCPREQUEST.

If the scope is configured with a domain name and `addZoneInHostname` is `true`, the domain name is appended to the client-provided hostname to form a fully qualified domain name (FQDN), as described [here](../scopes/#dns).

- Example with scope default options: `somehost`
- Example when scope has `addZoneInHostname` enabled: `somehost.example.com`

### IP Address Lease Time

- Tag: `51`

This option defaults to the TTL configured for the scope.
