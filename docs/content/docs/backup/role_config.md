---
title: "Role configuration"
---

- `accessKey`: S3 access key for backups.
- `secretKey`: S3 secret key for backups.

  If credentials aren't set, Gravity will try to use environment variables, and fall back to getting credentials from the EC2 service.

- `endpoint`: S3 endpoint.
- `bucket`: S3 bucket.
- `path`: S3 path prefix.
- `cronExpr`: Cron expression for backup frequency (defaults to `0 */24 * * *`).
