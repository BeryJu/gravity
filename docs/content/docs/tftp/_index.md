---
title: "TFTP Role"
---

Gravity provides a TFTP server that can store router/switch configurations and can be used for chain network booting.

## Reading

### Database

By default files are read from the Gravity database. Files are organized by the host's IP to prevent files from being overwritten by other hosts. To explicitly read files that should be accessible by other hosts, prefix the path with `shared/`.

### Bundled files

Gravity includes bundled files that can be used for TFTP network booting. These files can be accessed via the `bundled/` prefix.

Included files are:

 - `bundled/netboot.xyz.efi`: (UEFI) [Netboot](https://netboot.xyz) DHCP boot image file, uses built-in iPXE NIC drivers
 - `bundled/netboot.xyz.kpxe`: (Legacy) [Netboot](https://netboot.xyz) DHCP boot image file, uses built-in iPXE NIC drivers
 - `bundled/netboot.xyz-undionly.kpxe`: (Legacy) [Netboot](https://netboot.xyz) DHCP boot image file, use if you have NIC issues
 - `bundled/ipxe.undionly.kpxe`: (Legacy) [iPXE](https://ipxe.org) Chain image file

### Local files

(This option needs to be enabled in the Role settings)

Files stored in the `data/tftp` directory on individual nodes can be accessed via the `local/` prefix.

For example to download a local file `data/tftp/foo.bar`, enter the TFTP path `local/foo.bar`.

## Writing

### Database

Files are written to the Gravity database. Files are organized by the host's IP to prevent files from being overwritten by other hosts. To explicitly save files that should be accessible by other hosts, prefix the path with `shared/`.

Files have a maximum size of 10 MB.
