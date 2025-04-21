---
title: "external-dns Integrations"
weight: 0
---

## external-dns


### Configuration options

- `GRAVITY_URL`: URL to the Gravity instance this external-dns provider talks to.
- `GRAVITY_TOKEN`: Token used to authenticate to the gravity instance.
- `DOMAIN_FILTER`: Pipe-separated list of domains that external-dns should manage.

#### Common options

- `LOG_LEVEL`: Log level. Defaults to `info`.
- `DEBUG`: Enable debug mode. This should not be set manually in most cases and is only intended for development environments.
