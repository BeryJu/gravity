---
title: "Discovery Role"
weight: 12
description: Setup Gravity to discover existing devices on a network and create DNS/DHCP records/leases for them
---

Gravity can run periodical scans of your network(s) to discover devices. Once Gravity has discovered a device, you can create a DHCP lease and/or DNS record based on the discovered information.

### Concepts

---

##### Subnets

Each subnet is a layer 3 network. For each subnet, Gravity initiates a `nmap` ping scan to discover devices.

A subnet defines a CIDR to be scanned, a TTL indicating how long discovered devices should be saved for, and an optional DNS resolver for the subnet.

##### Devices

Each discovery is saved as a device. A device contains a MAC address, an IP address, and, if found, a hostname.

Devices can be "applied" to initiate the creation of either a DHCP lease or a DNS record (or possibly both, see below) based on the discovered information.

When you opt to create a DHCP lease while applying a device, a corresponding DNS record will also be generated, provided that the DHCP scope has the DNS integration configured (see [Scopes](../dhcp/scopes/#dns)).
