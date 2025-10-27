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

The subnet for which this scope is responsible.

#### `default`

Use this subnet as a fallback scope when no match can be found.

Default: `false`

#### `options`

List of dictionaries to set DHCP options.

Example options:

```yaml
- tagName: router
  value: 10.1.2.3
```

More info and default values can be found [here](./options).

#### `ttl`

TTL in seconds for the leases created from the scope.

Default: `86400`

#### `ipam`

Key:value settings for the IPAM

##### Internal IPAM (type = internal)

- `range_start`: Start of the range to give IPs from
- `range_end`: End of the range to give IPs from
- `should_ping`: Set to `"true"` to make Gravity ping an IP before giving it out.

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

#### `hook`

Optional hooks to dynamically modify requests and responses. See [Hooks](./hooks)
