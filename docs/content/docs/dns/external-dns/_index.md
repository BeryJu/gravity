---
title: "Integration with external-dns"
---

Webhook provider for [external-dns](https://kubernetes-sigs.github.io/external-dns/latest/) which allows external-dns to create DNS records in Gravity.

{{% alert title="Warning" color="warning" %}}
The gravity external-dns provider is in alpha.
{{% /alert %}}

## Example Helm chart value

```yaml
provider:
  name: webhook
  webhook:
    image:
      repository: ghcr.io/beryju/gravity-external-dns
      tag: stable
    securityContext:
      runAsNonRoot: true
      runAsUser: 1000
      runAsGroup: 1000
      allowPrivilegeEscalation: false
    env:
      - name: GRAVITY_URL
        value: http://my.gravity.instance:8008/
      - name: GRAVITY_TOKEN
        valueFrom:
          secretKeyRef:
            name: external-dns-secrets
            key: gravity-token
```

## Configuration options

- `GRAVITY_URL`: URL to the Gravity instance this external-dns provider talks to.
- `GRAVITY_TOKEN`: Token used to authenticate to the gravity instance.
- `DOMAIN_FILTER`: Pipe-separated list of domains that external-dns should manage.

### Common options

- `LOG_LEVEL`: Log level. Defaults to `info`.
- `DEBUG`: Enable debug mode. This should not be set manually in most cases and is only intended for development environments.
