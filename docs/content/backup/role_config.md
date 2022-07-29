---
title: "Role configuration"
---

- `accessKey`: S3 Access Key for backups
- `secretKey`: S3 Secret Key for backups

    If credentials aren't set, Gravity will try to use environment variables, and fall back to getting credentials from the EC2 service.

- `endpoint`: S3 Endpoint
- `bucket`: S3 bucket
- `path`: S3 Path prefix
- `cronExpr`: Cron expression on which backups are run (defaults to `0 */24 * * *`)
