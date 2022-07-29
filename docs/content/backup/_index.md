---
title: "Backup"
weight: 10
---

Gravity's inbuild backup regularly saves etcd snapshots to S3. These snapshots can then be used to restore a cluster from scratch.

Additionally, there are CLI tools to export/import all data in a more human-readable format:

### `cli export`

Running `gravity cli export` in the container will export all keys in the database into a JSON file.

Example:

```
gravity cli export /data/dump.json
```

### `cli import`

Running `gravity cli import` in the container will import all keys from a JSON file into the database.

Example:

```
gravity cli import /data/dump.json
```
