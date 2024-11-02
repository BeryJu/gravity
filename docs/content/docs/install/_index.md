---
title: "Installation"
weight: 5
description:
  Get started using Gravity
---

### Pre-requisites

Gravity requires at least 1 CPU core and 1 GB of memory. The resource usage varies depending on your configuration, for example having 200 DNS Zones with Blocky and CoreDNS will increase these requirements.

Gravity supports both x86 and ARM-based hosts, however only AMD64 and ARM64 variants are supported. Gravity should be able to work on other architectures too, however this is not officially supported.

### Docker Compose

Create a file called `docker-compose.yml` in a new directory with the following content:

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

Run `docker compose up -d` to start gravity.

## First time use

A default admin user is created on the first startup. You can find the credentials printed to stdout (accessible with `docker compose logs`).

You can reach Gravity by going to `http://<server IP or hostname>:8008` in your browser.

## Configuration

The following environment variables can be set.

### Common

- `BOOTSTRAP_ROLES`: Configure which roles this instance should bootstrap. Defaults to `dns;dhcp;api;etcd;discovery;backup;monitoring;tsdb;tftp`.
- `LOG_LEVEL`: Log level. Defaults to `info`.
- `DATA_PATH`: Path to store etcd data. Defaults to `./data`.
- `INSTANCE_IDENTIFIER`: Unique identifier of an instance, should ideally not change. Defaults to the detected hostname. When running in Docker, this is configured via the `hostname` attribute.
- `INSTANCE_IP`: This instance's reachable IP. When running in Docker and not using `network_mode: host`, this should be the host's IP.
- `LISTEN_ONLY`: Enable listen-only mode. In listen-only mode, Gravity will not reply to any DHCP packets and will not run [discovery](../discovery).

### Advanced

- `DEBUG`: Enable debug mode. This should not be set manually in most cases and is only intended for development environments.
- `INSTANCE_LISTEN`: By default the instance will listen on `INSTANCE_IP`, but this option will override that. Set to 0.0.0.0 when using Docker.
- `ADMIN_PASSWORD`: Optionally set a default password for the admin user. If unset, a random password will be generated as described [above](#first-time-use).
- `ADMIN_TOKEN`: Optionally set a token to be created on first start. If unset, no token will be created.
- `SENTRY_ENABLED`: Enable Sentry error reporting and tracing. Defaults to `false`.
- `SENTRY_DSN`: Configure a custom Sentry DSN.
- `ETCD_PREFIX`: Global etcd prefix. Defaults to `/gravity`.
- `ETCD_PEER_PORT`: Port used for etcd peer traffic. Defaults to `2380`. This may need to be changed when running in Kubernetes.
- `ETCD_ENDPOINT`: etcd Client endpoint. Defaults to `localhost:2379` when using embedded etcd.
- `ETCD_JOIN_CLUSTER`: Set to a join cluster token to join the node to a cluster. See [Clustering](./cluster).
- `IMPORT_CONFIGS`: A list of base64-encoded configs or file URIs to import configs from on first startup. For example:

    ```bash
    IMPORT_CONFIGS=file:///config.json
    # Multiple configs
    IMPORT_CONFIGS=file:///configA.json|file:///configB.json
    # Base64 import
    IMPORT_CONFIGS=eyJlbnRyaWVzIjogW3sia2V5IjogIi9ncmF2aXR5L2ZvbyIsInZhbHVlIjogIlptOXYifV19
    ```

### Changing Environment Variables

Gravity is designed so that you ideally don't have to explicitly define environment variables, but if necessary, environment variables can be added or changed by modifying the `environment:` options in the Docker Compose file.

Example:
```yaml
    environment:
      INSTANCE_IP: 192.168.2.8
      BOOTSTRAP_ROLES: dns;api;etcd;discovery;monitoring;tsdb
      INSTANCE_IDENTIFIER: my-gravity-server
```
