---
title: "From MS DHCP"
---

Gravity can import DHCP leases from Microsoft's DHCP server.

### Export leases and reservations

Run the following command on the server running Microsoft DHCP, or another server with management tools installed:

```powershell
Export-DhcpServer -Leases -File C:\dhcp_export.xml
```

This file contains information about all scopes, their leases and reservations, and any extra options.

Transfer the resulting .xml file to the server running gravity, and make sure it's available to gravity (it can be moved into the data container mount).

### Import data into gravity

Within the gravity container, run the following command to import the data

```bash
gravity cli convert ms_dhcp /data/dhcp_export.xml
```

This command will create all the required scopes, and convert leases and reservations into gravity equivalents.

Note that this does not configure DNS, so this has to either be done beforehand or afterwards.
