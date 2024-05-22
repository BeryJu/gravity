---
title: "API"
---

Gravity's API is accessible by default on port 8008 under `/api/v1`. OpenAPI is used to document the schema and automate the generation of client libraries.

Gravity's API is also available on the socket `/var/run/gravity.sock` with no authentication. This API is used by the `gravity cli` commands automatically.

### Authentication

Session authentication using local users and OIDC is supported for browser usage. API keys can also be created for automation.

A default admin user is created on the first startup. You can find the credentials printed to stdout. See [Installation](../install).

#### CLI

To create a new user, run the following command in the Gravity container:

```
gravity cli users add myusername
```

This will prompt you for a password which will be hashed and stored in the database.

The above command can also be used to reset a users' password, as it will overwrite any data for the given username.
