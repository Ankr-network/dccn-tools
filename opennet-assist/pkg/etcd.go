package pkg

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/pkg/transport"
)

func NewETCDClient(crt, key, ca string, endpoints []string) (*clientv3.Client, error) {
	log.Info().Msgf("ca: %s crt: %s key: %s endpoints: %+v", ca, crt, key, endpoints)
	tlsInfo := transport.TLSInfo{
		CertFile:      crt,
		KeyFile:       key,
		TrustedCAFile: ca,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 30 * time.Second,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return cli, nil
}
