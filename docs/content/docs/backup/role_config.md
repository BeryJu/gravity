---
title: "Role configuration"
---

## Local and S3 related settings

- Cron Schedule (`cronExpr`): Cron expression for backup frequency (defaults to `0 */24 * * *`, once every 24 hours).

## S3 related settings

Gravity only uploads files to the S3 bucket and as such only needs the `PutObject` permission.

- Bucket (`bucket`): S3 bucket.
- Path (`path`): S3 path prefix.
- Endpoint (`endpoint`): S3 endpoint.
- Access key (`accessKey`): S3 access key for backups.
- Secret key (`secretKey`): S3 secret key for backups.

  If credentials aren't set, Gravity will try to use environment variables, and fall back to getting credentials from the EC2 service.
