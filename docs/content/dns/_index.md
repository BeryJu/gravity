+++
title = "DNS"
date = 2022-07-17T21:08:57+02:00
weight = 6
pre = "<b>X. </b>"
+++

# DNS

Gravity's DNS Server supports

- resolving static hosts defined in the etcd database
- forwarding requests to other DNS Servers
- caching queries and responses in memory
- use [Blocky](https://0xerr0r.github.io/blocky/) for Ad and Privacy blocking
- use [k8s_gateway](https://github.com/ori-edge/k8s_gateway) to resolve kubernetes ingresses

### Concepts
---

##### Zones

Each DNS record belongs to a zone. Most commonly, the zone will be the domain part of an FQDN, so for the record `foo.bar.baz.`, it would be `bar.baz.`. Keep the trailing period in mind, as this is crucial for the zone to work properly.

Zones can also have lower level records, so for the zone `baz.`, you could add a record `foo.bar.` to get the same result as above.

The root zone, which is a zone for `.`, is used as fallback for any records for which a matching zone could not be found.

Each zone has it's individual configuration for how to handle queries, see [Handlers](./handlers) for more.

##### Records
