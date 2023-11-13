---
title: "Role configuration"
---

- `port`: The port the API server listens on (defaults to 8008).
- `oidc`: Optional OpenID Connect config.

  - `clientID`: OpenID Client Identifier.
  - `clientSecret`: OpenID Client Secret.
  - `issuer`: OpenID Issuer, sometimes also called "Configuration URL". Ensure `.well-known/openid-configuration` suffix is removed.
  - `redirectURL`: Redirect URL Gravity is reachable under. Should end in `/auth/oidc/callback`.

    The placeholder `$INSTANCE_IDENTIFIER` will be replaced by the instance's name and `$INSTANCE_IP` will be replaced by the instances IP.

  - `scopes`: Array of scopes that are requested. Should contain `openid` and `email`.
