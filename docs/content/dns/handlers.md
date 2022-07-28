+++
title = "Handlers"
+++

# Zone Handlers

The order of handler matters, gravity will send the query to each handler in the order the are configured, until one of the handlers returns a response.

### etcd

Attempt to reply to query by looking for records in etcd. Keep in mind because this is configured on zone level, this handler will only look for matching records in the current zone.

##### Configuration

None

### memory

Gravity watches the etcd database for changes, and keeps all the records in memory. If this handler is enabled for a zone which has records in etcd, this handler will reply from memory.

##### Configuration

None

### forward_ip

Forward queries to another DNS Server

##### Configuration

- `to` (required): List of DNS Servers to forward query to

    Multiple servers should be separated by `;`. For example `8.8.8.8:53;1.1.1.1`

- `cache_ttl`: Optional TTL to cache responses in etcd

    Defaults to 0. Attempts to cache for the TTL of the response.
    Set to -1 to never cache, and set to -2 to cache without a TTL

### forward_blocky

Forward queries to another DNS Server using blocky for Ad/Privacy blocking

##### Configuration

- `to` (required): List of DNS Servers to forward query to

    Multiple servers should be separated by `;`. For example `8.8.8.8:53;1.1.1.1`

- `cache_ttl`: Optional TTL to cache responses in etcd

    Defaults to 0. Attempts to cache for the TTL of the response.
    Set to -1 to never cache, and set to -2 to cache without a TTL

- `blocklists`: List of blocklists to load

    Multiple URLs should be separated by `;`. Defaults to these lists:

    - https://adaway.org/hosts.txt
    - https://dbl.oisd.nl/
    - https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext
    - https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
    - https://v.firebog.net/hosts/AdguardDNS.txt
    - https://v.firebog.net/hosts/Easylist.txt
    - https://adguardteam.github.io/AdGuardSDNSFilter/Filters/filter.txt

### k8s_gateway

Resolve queries by looking up Kubernetes Ingresses and services, see [k8s_gateway](https://github.com/ori-edge/k8s_gateway) for more info.

##### Configuration

None
