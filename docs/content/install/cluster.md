---
title: "Clustering"
---

Any number of Gravity instances can be clustered together. This has the following advantages:

- High availability of DNS
- Failover
- Multi-site with a central source of truth

Because gravity is based on etcd, it is recommended to use an odd number of nodes (see [etcd](https://etcd.io/docs/v3.5/faq/#why-an-odd-number-of-cluster-members)).
With 2 nodes in a cluster, you'll have the same quorum as with one node. Hence with 2 nodes, if any node goes down, the other node cannot access etcd either.
Gravity works around this issue with caching objects in memory, and as such even when etcd is not accessible, DNS resolution still works (as long as either forward_* or memory handlers are enabled).

## Adding nodes to a cluster

The recommended way to add a node is via the API `/api/v1/etcd/join` endpoint. This endpoint takes a `peer` parameter in the body of the request, which should contain the URL to the instance you're adding.

For example, given these two instances

- gravity-a, running on 10.100.1.10
- gravity-b, stopped on 10.100.2.10

The `peer` parameter would be set to `https://10.100.2.10:2380`.
