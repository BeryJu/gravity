---
title: Options
weight: 5
---

## Tag

### `tagName`:

The `tagName` corresponds to a human-readable name for a DHCP option tag; any of these tag names can be used:

- `subnet_mask`
- `router`
- `time_server`
- `name_server`
- `domain_name`
- `bootfile`
- `tftp_server`

*Conflicts with `tag`*

### `tag`:

This field allows for setting the raw DHCP option tag to any value.

*Conflicts with `tagName`*

## Value

### `value`:

Set the value this option should be set to. When used in conjunction with `tagName`, Gravity will automatically encode the value correctly. When used with `tag`, the value needs to be correctly escaped manually.

Example:

```yaml
- tagName: router
  value: 10.1.2.3
```

*Conflicts with `value64` and `valueHex`*

### `value64`:

Set the value this option should be set to, encoded in base64. This allows for pre-encoding the value when using `tag`, and representing data that can't be encoded in an ASCII string. This value should be set to an array of base64-strings, all of which are concatenated.

*Conflicts with `value`*

### `valueHex`:

Set the value this option should be set to, encoded in Hexadecimal. This value should be set to an array of hex-strings, all of which are concatenated.

Example:

```yaml
- tag: 43
  valueHex:
    - 0104C0A8030A
```

*Conflicts with `value`*
