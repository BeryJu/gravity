---
title: "Zones"
weight: 1
---

## Handlers

The order of handler matters; Gravity will send a query to each handler in the order they are configured until a response is returned.

The handler configuration consists of a list of individual handler configurations. All list entries require a `type` attribute which must match one of the headers listed below. For example:

```yaml
- cache_ttl: 3600
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

Attempt to reply to query by looking for records in etcd. Keep in mind that because this is configured on zone level, this handler will only look for matching records in the current zone.

##### Configuration

None.

### `memory`

Gravity watches the etcd database for changes and keeps all the records in memory. If this handler is enabled for a zone which has records in etcd, this handler will reply from memory.

##### Configuration

None.

### `forward_ip`

Forward queries to another DNS server.

##### Configuration

- `to` (required): List of DNS Servers to forward query to.

  Multiple servers should be separated by `;`. For example `8.8.8.8:53;1.1.1.1`.

  Starting with Gravity 0.14, this can also be set as a JSON/YAML array using `[8.8.8.8:53, 1.1.1.1:53]` instead of a semicolon-separated string.

- `cache_ttl`: Optional TTL to cache responses in etcd.

  Defaults to 0. Attempts to cache for the TTL of the response.
  Set to -1 to never cache, and set to -2 to cache without a TTL.

##### Example

```yaml
- type: forward_ip
  to: 8.8.8.8
```

### `forward_blocky`

Forward queries to another DNS server via Blocky for advert/privacy blocking.

##### Configuration

- `to` (required): List of DNS Servers to forward query to (via Blocky).

  Multiple servers should be separated by `;`. For example `8.8.8.8:53;1.1.1.1`.

  Starting with Gravity 0.14, this can also be set as a JSON/YAML array using `[8.8.8.8:53, 1.1.1.1:53]` instead of a semicolon-separated string.

- `cache_ttl`: Optional TTL to cache responses in etcd

  Defaults to 0. Attempts to cache for the TTL of the response.
  Set to -1 to never cache, and set to -2 to cache without a TTL.

- `config`: Optional Blocky configuration as string. (Requires Gravity 0.15.0)

  See [here](https://0xerr0r.github.io/blocky/main/configuration/) for a reference configuration file and options that can be configured.

- `blocklists`: List of blocklists to load.
- `allowlists`: List of allowlists to load.

  Entries beginning with http:// or https:// are downloaded when the DNS Role is started. Can also be set to an inline list of domains to block/allow. Multiple entries should be separated by a `;`.

  Starting with Gravity 0.14, this can also be set as a JSON/YAML array using `[https://foo, https://bar]` instead of a semicolon-separated string.

  By default, these blocklists are loaded:

  - https://adaway.org/hosts.txt
  - https://dbl.oisd.nl/
  - https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
  - https://v.firebog.net/hosts/AdguardDNS.txt
  - https://v.firebog.net/hosts/Easylist.txt
  - https://adguardteam.github.io/AdGuardSDNSFilter/Filters/filter.txt

  These default lists are downloaded at compile time and embedded into Gravity in order to speed up the startup of Blocky.

##### Example

{{< tabpane text=true >}}
  {{% tab header="Customized options" lang="en" %}}

```yaml
- type: forward_blocky
  to: 8.8.8.8
  blocklists:
    - https://adaway.org/hosts.txt
  allowlists:
    - exception.com
```

  {{% /tab %}}
  {{% tab header="Full custom config" lang="en" %}}

```yaml
- type: forward_blocky
  # Non-exhaustive Blocky configuration, this is just an example to show the usage of `config:`
  config: |
    upstreams:
      init:
        # Configure startup behavior.
        # accepted: blocking, failOnError, fast
        # default: blocking
        strategy: fast
      groups:
        # these external DNS resolvers will be used. Blocky picks 2 random resolvers from the list for each query
        # format for resolver: [net:]host:[port][/path]. net could be empty (default, shortcut for tcp+udp), tcp+udp, tcp, udp, tcp-tls or https (DoH). If port is empty, default port will be used (53 for udp and tcp, 853 for tcp-tls, 443 for https (Doh))
        # this configuration is mandatory, please define at least one external DNS resolver
        default:
          # example for tcp+udp IPv4 server (https://digitalcourage.de/)
          - 5.9.164.112
          # Cloudflare
          - 1.1.1.1
```

  {{% /tab %}}
{{< /tabpane >}}

### `coredns`

Resolve queries by using a variety of CoreDNS Plugins. See [here](https://coredns.io/plugins/) for all plugins.

##### Configuration

- `config`: String configuration in the caddyfile format that is passed to CoreDNS.

  Example:

  ```caddy
  .:1053 {
    whoami
  }
  ```

  **Make sure your configuration uses a port other than the one Gravity uses to prevent any issues.**

##### Example

```yaml
- type: coredns
  config: |
    .:1053 {
      whoami
    }
```

## Hooks

Optional hooks to dynamically modify requests and responses. See [Hooks](../hooks)
