---
title: "Migration as a proxy"
linkTitle: "as a proxy"
---

A universal method of migrating DNS records to Gravity is to setup a Zone to forward requests to an existing DNS Server and save the results into Gravity.

This can be done by using the following handler configuration on a zone:

```yaml
- type: memory
- type: etcd
- type: forward_ip
  to:
    - <IP Address of the existing DNS Server(s)>
  cache_ttl: -2
```

With this configuration, Gravity can be configured as the default DNS Server for a domain, and it will forward requests to the configured upstream DNS servers and permanently write the results into its database.
