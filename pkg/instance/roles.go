package instance

import (
	_ "beryju.io/gravity/pkg/roles/api"
	_ "beryju.io/gravity/pkg/roles/backup"
	_ "beryju.io/gravity/pkg/roles/debug"
	_ "beryju.io/gravity/pkg/roles/dhcp"
	_ "beryju.io/gravity/pkg/roles/discovery"
	_ "beryju.io/gravity/pkg/roles/dns"
	_ "beryju.io/gravity/pkg/roles/etcd"
	_ "beryju.io/gravity/pkg/roles/monitoring"
	_ "beryju.io/gravity/pkg/roles/tftp"
	_ "beryju.io/gravity/pkg/roles/tsdb"
)
