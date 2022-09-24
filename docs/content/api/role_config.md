---
title: "Role configuration"
---

- `port`: The Port the API server listens on (defaults to 8008)
- `oidc`: Optional OpenID Connect config

    - `clientID`: OpenID Client Identifier
    - `clientSecret`: OpenID Client Secret
    - `issuer`: OpenID Issuer, sometimes also called "Configuration URL". Ensure `.well-known/openid-configuration` suffix is removed.
    - `redirectURL`: Redirect URL Gravity is reachable under, should end in `/auth/oidc/callback`.
    - `scopes`: Array of scopes that are requested, should contain `openid` and `email`
