---
title: "Installation"
weight: 5
---

### Installation

#### docker-compose

```yaml
---
version: "3.4"

services:
  gravity:
    # Important for this to be static and unique
    hostname: gravity1
    image: ghcr.io/beryju/gravity:stable
    restart: unless-stopped
    network_mode: host
    volumes:
      - data:/data
    # environment:
    #   LOG_LEVEL: info
    # The default log level of info logs DHCP and DNS queries, so ensure
    # the logs aren't filling up the disk
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  data:
    driver: local
```

## First time use

A default admin user is created on the first startup. You can find the credentials printed to stdout.

### Configuration

The following environment variables can be set

##### Common

- `BOOTSTRAP_ROLES`: Configure while roles this instance should bootstrap, defaults to `dns;dhcp;api;etcd;discovery;backup;monitoring;tsdb`.
- `LOG_LEVEL`: Log level, defaults to `info`.
- `DATA_PATH`: Path to store etcd data, defaults to `./data`
- `INSTANCE_IDENTIFIER`: Unique identifier of an instance, should ideally not change. Defaults to hostname. When running in docker, this is configured via the `hostname` attribute.
- `INSTANCE_IP`: This instance's reachable IP, when running in docker and not using `network_mode: host`, this should be the hosts IP
- `LISTEN_ONLY`: Enable listen-only mode which will not reply to any DHCP packets and not run discovery

##### Advanced

- `DEBUG`: Enable debug mode
- `ETCD_PREFIX`: Global etcd prefix, defaults to `/gravity`
- `ETCD_ENDPOINT`: etcd Client endpoint, defaults to `localhost:2379` when using embedded etcd
- `ETCD_JOIN_CLUSTER`: Used when joining a node to a cluster, value is given by join API endpoint
- `INSTANCE_LISTEN`: By default the instance will listen on `INSTANCE_IP`, but can be set to override that (set to 0.0.0.0 in docker)
- `ADMIN_PASSWORD`: Optionally set a default password for the admin user, if not set a random one will be generated
- `ADMIN_TOKEN`: Optionally set a token to be created on first start, if not set no token will be created
- `SENTRY_ENABLED`: Enable sentry error reporting and tracing
- `SENTRY_DSN`: Configure a custom sentry DSN
