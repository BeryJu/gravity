---
title: "DHCP Role"
weight: 6
description: Configure Gravity as your DHCP server and optionally import existing leases/reservations.
---

Gravity's DHCP Server supports:

- Reservations
- Static leases
- Full replication
- Pluggable IPAMs

### Concepts

---

##### Scopes

A scope represents a single layer 2 network. Within a scope, Gravity assumes to be the main and only DHCP server, when not in listen-only mode.

The scope for any particular request is found in the following way:

- Checking if the DHCP Peer address is within any of the scopes CIDR, or
- Checking if the Address of the interface the request was received on is within any of the scopes CIDR, or
- Checking if the scope has the `default` flag set to `true`.

Afterwards, all matching scopes are sorted by the length of the scope's CIDR Prefix, and the longest prefix match is chosen.

See [Scopes](./scopes).

##### Leases

A lease represents a single IP address assigned to a single MAC Address. A lease also includes a hostname.
