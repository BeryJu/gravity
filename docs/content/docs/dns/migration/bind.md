---
title: "Migrating from Bind"
linkTitle: "from Bind"
---

(requires Gravity v0.24+)

Gravity can import DNS records from Bind zonefiles.

### Retrieving the zone file

Retrieving the zone file will depend on the setup of your Bind server. In most cases, this file will be in `/etc/bind/`.

Some other DNS servers also provide the option of exporting records as DNS, these files are also compatible with Gravity.

### Import data into Gravity

##### Web

Click on the **Create** button on the *DNS Zones* page and enter the domain for the new zone. On the next page select **Import** and continue. Select the zone file from above. Gravity will import all records into the new scope.

##### CLI

The zone file from above must be transferred to the server running Gravity, and must be accessible in the Gravity container (it can be moved into the `/data` container mount).

Within the Gravity container, run the following command to import the data:

```bash
gravity cli convert bind /data/mydoamin.zone
```

This command will create all the required zones, converting records into Gravity equivalents.
