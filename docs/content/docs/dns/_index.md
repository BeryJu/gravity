---
title: "DNS"
weight: 6
---

Gravity's DNS Server supports

- resolving static hosts defined in the etcd database
- forwarding requests to other DNS Servers
- caching queries and responses in memory
- use [Blocky](https://0xerr0r.github.io/blocky/) for Ad and Privacy blocking
- use [CoreDNS](https://github.com/coredns/coredns) to resolve use any CoreDNS plugin

### Concepts

---

##### Zones

Each DNS record belongs to a zone. Most commonly, the zone will be the domain part of an FQDN, so for the record `foo.bar.baz.`, it would be `bar.baz.`. Keep the trailing period in mind, as this is crucial for the zone to work properly.

Zones can also have lower level records, so for the zone `baz.`, you could add a record `foo.bar.` to get the same result as above. The longest matching zone is picked to resolve a record. If all of the handlers of a zone return no response, there is no fallthrough to the next zone.

The root zone, which is a zone for `.`, is used as fallback for any records for which a matching zone could not be found.

Each zone has it's individual configuration for how to handle queries, see [Handlers](./handlers) for more.

##### Records

A record belongs to one zone and stores one response. To support multiple responses (i.e. multiple IP addressess for an A record), Record UIDs are used. A UID is optional, and records with UID can be combined with a record without UID (all their results will be returned). Records created by the DHCP role will automatically have the UID assigned based on the DHCP devices identifier (the MAC address in most cases).

A single record holds the following data:

- `data`: The actual response, an IP for A/AAAA records, Text for TXT records, etc
- `ttl`: TTL of the response, optional

_For MX records_

- `mxPreference`: Configure the MX Preference (optional)

_For SRV records_

- `srvPort`: Configure SRV Port (optional)
- `srvPriority`: Configure SRV Priority (optional)
- `srvWeight`: Configure SRV Weight (optional)
