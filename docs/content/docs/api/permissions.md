---
title: "Permissions"
---

Starting with Gravity 0.16, users can have permissions assigned to them. Permissions are assigned based on HTTP URL paths and methods. For example, to give a user permissions to read all resources, permissions can be set to this:

```json
[
  {
    "path": "/*",
    "methods": ["get", "head"]
  }
]
```

To give users permissions to view DNS zones and records, you can set the permissions to

```json
[
  {
    "path": "/api/v1/dns/zones",
    "methods": ["get"]
  },
  {
    "path": "/api/v1/dns/zones/records",
    "methods": ["get"]
  },
]
```

To give a user admin permissions, set the permissions to

```json
[
  {
    "path": "/*",
    "methods": ["*"]
  },
]
```
