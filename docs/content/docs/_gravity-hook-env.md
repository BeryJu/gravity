### `gravity` Object

#### `log(msg: any)`

Logs a message to the stdout of the Gravity node this hook is run on.

#### `node: string`

The identifier of the node this hook is run on.

#### `version: string`

The version of Gravity on the node this hook is run on.

#### `role`

A reference to the [Role instance](https://pkg.go.dev/beryju.io/gravity/pkg/instance#RoleInstance) this hook was triggered by.

### `net` Object

#### `parseIP(ip: string, family: string)`

Parse an IP address from the string `ip` and return it as an array of bytes. `family` determines if the IP should be parsed as IPv4 or IPv6.</p>
