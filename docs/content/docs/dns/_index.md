---
title: "DNS"
weight: 6
---

Gravity's DNS Server supports:

- Resolving static hosts defined in the etcd database
- Forwarding requests to other DNS Servers
- Caching queries and responses in memory
- Using [Blocky](https://0xerr0r.github.io/blocky/) for DNS-based advert/privacy blocking
- Using [CoreDNS](https://github.com/coredns/coredns) plugins

### Concepts

---

##### Zones

Each DNS record belongs to a zone. Most commonly, the zone will be the domain part of an FQDN, so for the record `foo.bar.baz.`, it would be `bar.baz.`. Keep the trailing period in mind, as this is crucial for the zone to work properly.

Zones can also have lower level records, so for the zone `baz.`, you could add a record `foo.bar.` to get the same result as above. The longest matching zone is picked to resolve a record. If all the handlers for a zone return no response, there is no fall-through to the next zone.

The root zone, which is a zone for `.`, is used as fallback for any records for which a matching zone could not be found.

Each zone has its individual configuration for how to handle queries, see [Handlers](./handlers) for more.

##### Records

A record belongs to one zone and stores one response. To support multiple responses (i.e. multiple IP addresses for an A record), record UIDs are used. A UID is optional, and records with UID can be combined with a record without UID (all their results will be returned). Records created by the DHCP role will automatically have the UID assigned based on the DHCP device's identifier (the MAC address in most cases).

To create a record at the root of the zone, set the name of the record to `@`.

To create a wildcard record, set the name of the record to `*`. Note that if a more specific record exists for the queried name, it will have a higher priority and the wildcard record will not be returned.
Wildcard records can also be used for multiple levels, for example creating a record called `*.*` in a zone `example.com` will be matched for a query to `foo.bar.example.com`. Here the first wildcard record, sorted by the least amount of depth (amount of `.`) will be returned, and no other records will be returned.

A single record holds the following data:

- `data`: The actual response, an IP for A/AAAA records, text for TXT records, etc.
- `ttl`: TTL of the response (optional).

_For MX records_

- `mxPreference`: Configure the MX Preference (optional).

_For SRV records_

- `srvPort`: Configure SRV Port (optional).
- `srvPriority`: Configure SRV Priority (optional).
- `srvWeight`: Configure SRV Weight (optional).
