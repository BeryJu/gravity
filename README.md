# gravity

Fully-replicated DNS and DHCP Server using etcd

![](./docs/content/images/dns_zones.png)
![](./docs/content/images/dns_records.png)
![](./docs/content/images/dhcp_leases.png)
![](./docs/content/images/cluster_nodes.png)
![](./docs/content/images/grafana.png)

### Configuration

##### External Configuration

The following environment variables can be set:

- `DEBUG`: Enable debug mode
- `LOG_LEVEL`: Log level, defaults to `info`.
- `DATA_PATH`: Path to store etcd data, defaults to `./data`
- `ETCD_PREFIX`: Global etcd prefix, defaults to `/gravity`
- `ETCD_ENDPOINT`: etcd Client endpoint, defaults to `localhost:2379` when using embedded etcd
- `ETCD_JOIN_CLUSTER`: Used when joining a node to a cluster, value is given by join API endpoint
- `BOOTSTRAP_ROLES`: Configure while roles this instance should bootstrap, defaults to `dns;dhcp;api;etcd;discovery;backup`.
- `INSTANCE_IDENTIFIER`: Unique identifier of an instance, should ideally not change. Defaults to hostname, but needs to be set in containers.
- `INSTANCE_IP`: This instance's reachable IP, when running in docker this should be the hosts IP
- `INSTANCE_LISTEN`: By default the instance will listen on `INSTANCE_IP`, but can be set to override that (set to 0.0.0.0 in docker)
- `LISTEN_ONLY`: Enable listen-only mode which will not reply to any DHCP packets and not run discovery

### API

Checkout `http://$INSTANCE_IP:8008/api/v1/docs`
