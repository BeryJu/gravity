---
title: "Discovery"
---

Gravity can run periodical scans of your network(s) to discover devices. These discovered devices can then be turned into DHCP Leases or DNS Records.

### Concepts
---

##### Subnets

Each subnet is a Layer3 network that Gravity will scan. A subnet defines a CIDR which will be scanned, and a TTL for how long discovered devices should be saved for.

##### Devices

Each discovery is saved as a Device, which contains a MAC address, an IP Address and possibly a hostname.

Devices can be *applied* to convert them into DHCP Leases, which will also create DNS Records if the DHCP scope has the DNS integration eanbled, or just DNS records without DHCP.
