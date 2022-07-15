package etcd

import (
	"errors"
	"os"
	"path"
)

const (
	relCACertPath = "/ca.pem"
	relCAKeyPath  = "/ca_key.pem"

	relInstCertPath = "/instance.pem"
	relInstKeyPath  = "/instance_key.pem"
)

func (ee *Role) configureCertificates() {
	// Always check for CA public key, if we don't have one
	// assume new install and generate a CA
	if _, err := os.Stat(path.Join(ee.certDir, relCACertPath)); errors.Is(err, os.ErrNotExist) {
		err := ee.generateCA()
		if err != nil {
			ee.log.WithError(err).Warning("failed to generate CA")
			return
		}
	}
	// Check Instance CA, if it doesn't exist attempt to load
	// the CA and create one
	if _, err := os.Stat(path.Join(ee.certDir, relInstCertPath)); errors.Is(err, os.ErrNotExist) {
		err := ee.generateInstance()
		if err != nil {
			ee.log.WithError(err).Warning("failed to generate instance cert")
			return
		}
	}
	_, _, err := ee.loadInstance()
	if err != nil {
		ee.log.WithError(err).Warning("failed to load instance cert")
		return
	}
	ee.cfg.PeerTLSInfo.CertFile = path.Join(ee.certDir, relInstCertPath)
	ee.cfg.PeerTLSInfo.KeyFile = path.Join(ee.certDir, relInstKeyPath)
}
