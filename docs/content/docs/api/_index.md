---
title: "API Role"
weight: 10
description: Options related to Gravity's API and Web interface
---

Gravity's API is accessible by default on port 8008 under `/api/v1`. OpenAPI is used to document the schema and automate the generation of client libraries.

Gravity's API is also available on the socket `/var/run/gravity.sock` with no authentication. This API is used by the `gravity cli` commands automatically.

### Authentication

A default admin user is created on the first startup. You can find the credentials printed to stdout. See [Installation](../install).

Session authentication using local users and OIDC is supported for browser usage. API keys can also be created for automation.

To authenticate to the API using a token, create the token either using [ADMIN_TOKEN](../install/#advanced), or in the Web Interface under __Auth -> Tokens__. Upon creation, the token will be shown in the browser. Afterwards, add the `Authorization` header to API requests with the value of `Bearer <token>`.

Starting with Gravity 0.19, when OIDC is enabled, JWT tokens signed by the OIDC issuer can also be used. The role configuration parameter `tokenUsernameField` configures which claim from the JWT is used to lookup the user.

#### CLI

To create a new user, run the following command in the Gravity container:

```
gravity cli users add myusername
```

This will prompt you for a password which will be hashed and stored in the database.

The above command can also be used to reset a users' password, as it will overwrite any data for the given username.
