---
title: "Backup Role"
weight: 10
description: Configure Gravity's inbuilt backup system to upload snapshots to S3 or create them locally.
---

## Snapshots

Gravity's inbuilt backup regularly saves etcd snapshots to a user-configured S3 endpoint and bucket. These snapshots can then be used to restore a cluster from scratch.

Starting with 0.4.5, Gravity will also keep 1 snapshot locally without any additional configuration, in the `/data/backup` directory.
This snapshot is created/updated with the same schedule as configured in the [Role configuration's `Cron Schedule`](./role_config.md#local-and-s3-related-settings)
This snapshot can be restored just like other snapshots.

### Restore

1. Move the snapshot file into a volume accessible by the Gravity container.
2. In the Gravity container, run `gravity cli snapshot status /path/to/snapshot` to ensure the snapshot is a valid file.
3. Run `gravity cli snapshot restore /path/to/snapshot --data-dir /data/restore`. This will restore the snapshot to `/data/restore`, but the following steps are required to make Gravity use the restored data.
4. Stop the container.
5. Remove the `etcd` folder in the `data` volume.
6. Rename `restore` to `etcd`.
7. Start the container.
8. Run `gravity cli etcdctl member list -w table` in the Gravity container to list all nodes.
9. If the value under `PEER ADDRS` is `http://localhost:2380`, run the following command:

    ```
    # Replace <ID> with the value `ID` from the command above.
    gravity cli etcdctl member update <ID> --peer-urls http://$INSTANCE_IP:2380
    ```

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
