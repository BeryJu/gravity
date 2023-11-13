---
title: "Zone Handlers"
weight: 1
---

The order of handler matters, Gravity will send the query to each handler in the order the are configured, until one of the handlers returns a response.

The handler configuration consists of a list of individual handler configurations. All list entries require a `type` attributes, which must match with one of the headers listed below, for example:

```yaml
- cache_ttl: "3600"
  to: 8.8.8.8:53
  type: forward_blocky
- to: 8.8.8.8:53
  type: forward_ip
```

or

```yaml
- type: memory
- type: etcd
```

### `etcd`

Attempt to reply to query by looking for records in etcd. Keep in mind because this is configured on zone level, this handler will only look for matching records in the current zone.

##### Configuration

None

### `memory`

Gravity watches the etcd database for changes, and keeps all the records in memory. If this handler is enabled for a zone which has records in etcd, this handler will reply from memory.

##### Configuration

None

### `forward_ip`

Forward queries to another DNS Server

##### Configuration

- `to` (required): List of DNS Servers to forward query to

  Multiple servers should be separated by `;`. For example `8.8.8.8:53;1.1.1.1`

- `cache_ttl`: Optional TTL to cache responses in etcd

  Defaults to 0. Attempts to cache for the TTL of the response.
  Set to -1 to never cache, and set to -2 to cache without a TTL

##### Example

```yaml
- type: forward_ip
  to: 8.8.8.8
```

### `forward_blocky`

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
  - https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
  - https://v.firebog.net/hosts/AdguardDNS.txt
  - https://v.firebog.net/hosts/Easylist.txt
  - https://adguardteam.github.io/AdGuardSDNSFilter/Filters/filter.txt

  (These default lists are downloaded at compile time and embedded into Gravity, in order to speed up the startup of Blocky)

##### Example

```yaml
- type: forward_blocky
  to: 8.8.8.8
  blocklists: https://adaway.org/hosts.txt
```

### `coredns`

Resolve queries by using a variety of CoreDNS Plugins. See [here](https://coredns.io/plugins/) for all plugins.

##### Configuration

- `config`: String configuration in the caddyfile format that is passed to CoreDNS.

  Example:

  ```
  .:1053 {
    whoami
  }
  ```

  **Make sure to use a different port in the configuration as the one Gravity uses to prevent any issues**

##### Example

```yaml
- type: coredns
  config: |
    .:1053 {
      whoami
    }
```
