---
title: "Backup"
weight: 10
---

## Snapshots

Gravity's inbuilt backup regularly saves etcd snapshots to S3. These snapshots can then be used to restore a cluster from scratch.

Starting with 0.4.5, Gravity will also keep 1 snapshot locally without any additional configuration, in the `/data/backup` directory. This snapshot can be restored just like other snapshots

### Restore

Move the snapshot file into a volume accessible by the Gravity container.

In the Gravity container run `gravity cli snapshot status /path/to/snapshot`, which will ensure the snapshot is a valid file.

Run `gravity cli snapshot restore /path/to/snapshot --data-dir /data/restore`. This will restore the snapshot to `/data/restore`. To make Gravity use this data, stop the container, remove the `etcd` folder in the `data` volume, and rename `restore` to `etcd`. Afterwards start the container again.

Snapshots don't include any cluster information, so if you were using a cluster, you'll have to re-join nodes.

## CLI

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
