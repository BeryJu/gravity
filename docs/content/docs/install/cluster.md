---
title: "Clustering"
---

Any number of Gravity instances can be clustered together. This has the following advantages:

- High availability of DNS
- Failover
- Multi-site with a central source of truth

Because gravity is based on etcd, it is recommended to use an odd number of nodes (see [etcd](https://etcd.io/docs/v3.5/faq/#why-an-odd-number-of-cluster-members)).
With 2 nodes in a cluster, you'll have the same quorum as with one node. Hence with 2 nodes, if any node goes down, the other node cannot access etcd either.
Gravity works around this issue with caching objects in memory, and as such even when etcd is not accessible, DNS resolution still works (as long as either [`forward_*` or `memory` handlers](../../dns/handlers) are enabled).

## Adding nodes to a cluster

Adding a node can be done via the web interface, which will output a ready-to-use Docker Compose file.

Alternatively, a token can be created, which can be used to join nodes to the cluster. To join a node manually, set the `ETCD_JOIN_CLUSTER` environment variable.

The format should be the API token followed by a command followed by the HTTP URL of the instance that should be used to join off of.

For example, given the token `bootstrap-token` and the IP `10.0.0.1`, the value should be:

`bootstrap-token,http://10.0.0.1:8008`
