---
title: "Scopes"
weight: 1
---

The scope for any particular request is found in the following way:

- Checking if the DHCP Peer address is within any of the scopes CIDR, or
- Checking if the Address of the interface the request was received on is within any of the scopes CIDR, or
- Checking if the Scope has the `default` flag set to `true`.

Afterwards, all matching scopes are sorted by the length of the Scope's CIDR Prefix, and the longest prefix match is chosen.

## Config

#### `subnetCidr`

The CIDR this subnet is responsible for

#### `default`

Use this Subnet as a fallback scope when no match can be found

#### `options`

List of dictionaries to set DHCP options.

Example options:

```yaml
- tagName: router
  value: 10.1.2.3
```

More info [here](./scopes.md).

#### `ttl`

TTL for the leases created from the scope

#### `ipam`

Key:value settings for the IPAM

##### Internal IPAM (type = internal)

- `range_start`: Start of the range to give IPs from
- `range_end`: End of the range to give IPs from
- `should_ping`: Set to `"true"` to make gravity ping an IP before giving it out.

#### `dns`

DNS settings for this scope. If the zone exists as a zone in Gravity, then DNS integration is enabled. This will automatically create forward and reverse records with the same TTL as the lease (the records will also be renewed).

Additionally, `addZoneInHostname` can be set to make Gravity append the zone to the DHCP hostname

```json
{
  "dns": {
    "zone": "foo.bar.baz.",
    "search": [],
    "addZoneInHostname": false
  }
}
```
