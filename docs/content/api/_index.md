---
title: "API"
---

Gravity's API is accessible by default on Port 8008 under `/api/v1`. OpenAPI is used to document the schema and automate the generation of client libraries.

### Authentication

Session authentication using local users and OIDC is supported for browser usage. API Keys can also be created for automation.

A default admin user is created on the first startup. You can find the credentials printed to stdout.

#### CLI

To create a new user, run the following command in the gravity container:

```
gravity cli users add myusername
```

This will prompt you for a password, which will be hashed and stored in the database.
