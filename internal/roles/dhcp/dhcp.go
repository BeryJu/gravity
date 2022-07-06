package dhcp

import (
	"net"

	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles"
	log "github.com/sirupsen/logrus"

	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"github.com/insomniacslk/dhcp/dhcpv6/server6"
)

const (
	RoleDHCPPrefix = "dhcp"
)

type DHCPServerRole struct {
	s4  *server4.Server
	s6  *server6.Server
	log *log.Entry
}

func New(instance roles.Instance) *DHCPServerRole {
	return &DHCPServerRole{
		log: log.WithField("role", "dhcp"),
	}
}

func (r *DHCPServerRole) startServer4() error {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 67,
	}
	server, err := server4.NewServer("", &laddr, r.handler4)
	if err != nil {
		return err
	}
	r.s4 = server
	return r.s4.Serve()
}

func (r *DHCPServerRole) Start(config []byte) error {
	return nil
	// return r.startServer4()
}

func (r *DHCPServerRole) Stop() {
	r.s4.Close()
}

func (r *DHCPServerRole) HandleEvent(ev *roles.Event[any]) {
}
