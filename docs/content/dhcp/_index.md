---
title: "DHCP"
weight: 6
---

Gravity's DHCP Server supports

- reservations
- static leases
- full replication
- pluggable IPAMs

### Concepts
---

##### Scopes

A scope represents a single Layer-2 network. Within a scope, gravity assumes to be the main and only DHCP server, when not in listen-only mode.

The scope for any particular request is found in the following way:

- Checking if the DHCP Peer address is within any of the scopes CIDR, or
- Checking if the Address of the interface the request was received on is within any of the scopes CIDR, or
- Checking if the Scope has the `default` flag set to `true`.

Afterwards, all matching scopes are sorted by the length of the Scope's CIDR Prefix, and the longest prefix match is chosen.

See [Scopes](./scopes)

##### Leases

A lease represents a single IP Address assigned to a single MAC Address. A lease also includes a hostname
