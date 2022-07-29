---
title: "API"
---

Gravity's API is accessible by default on Port 8008 under `/api/v1`. OpenAPI is used to document the schema and automate the generation of client libraries.

### Authentication

Currently, only basic authentication is supported. Users can be added via the API or CLI.

#### CLI

To create a new user, run the following command in the gravity container:

```
gravity cli addUser -u myusername
```

This will prompt you for a password, which will be hashed and stored in the database.
