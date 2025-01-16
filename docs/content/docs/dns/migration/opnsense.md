---
title: "Migrating from Opensense"
linkTitle: "from Opensense"
---

(requires Gravity v0.24+)

Gravity can import DNS records from Opnsense backups.

### Exporting records

In the Opnsense web interface, export a backup under **System -> Configuration -> Backups**. Make sure to not include RRD data (default) and not encrypt the configuration file.

### Import data into Gravity

##### Web

Click on the **Create** button on the *DNS Zones* page and enter the domain for the new zone. On the next page select **Import** and continue. Select the zone file from above. Gravity will import all records into the new scope.

##### CLI

The resulting .xml file from above must be transferred to the server running Gravity, and must be accessible in the Gravity container (it can be moved into the `/data` container mount).

Within the Gravity container, run the following command to import the data:

```bash
gravity cli convert opnsense /data/opnsense.xml
```

This command will create all the required zones, converting records into Gravity equivalents.
