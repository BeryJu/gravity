---
title: "Migrating from Microsoft DHCP"
linkTitle: "from Microsoft DHCP"
---

Gravity can import DHCP leases from Microsoft's DHCP server.

### Export leases and reservations

Run the following command on the server running Microsoft DHCP, or another server with management tools installed:

```powershell
Export-DhcpServer -Leases -File C:\dhcp_export.xml
```

This file contains information about all scopes, their leases and reservations, and any extra options.

### Import data into Gravity

##### Web (requires Gravity v0.24+)

Click on the **Create** button on the *DHCP Scopes* page and enter the name for a new scope. On the next page select **Import** and continue. Select the XML file exported above. Gravity will import all leases and reservations into the new scope.

##### CLI

The resulting .xml file from above must be transferred to the server running Gravity, and must be accessible in the Gravity container (it can be moved into the `/data` container mount).

Within the Gravity container, run the following command to import the data:

```bash
gravity cli convert ms_dhcp /data/dhcp_export.xml
```

This command will create all the required scopes, converting leases and reservations into Gravity equivalents.

Note that this does not configure DNS, so this has to either be done beforehand or afterwards.
