package nats

import (
	"crypto/tls"

	"github.com/BloodRedCodr/go-utils/auth"
	"github.com/BloodRedCodr/go-utils/logger"
	"github.com/nats-io/nats.go"
)

func Connect(
	logger *logger.Logger,
	p12CertPath string, p12Pwd string,
	host string, port string,
) *nats.Conn {
	cert, caCertPool, _, _ := auth.GetCertsFromP12(logger, p12CertPath, p12Pwd)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   host,
		MinVersion:   tls.VersionTLS12,
	}

	natsServer := "nats://" + host + ":" + port
	nc, err := nats.Connect(natsServer, nats.Secure(tlsConfig))
	if err != nil {
		logger.Fatal("Error connecting to NATS: %v", err)
	}
	logger.Info("NATS connection established successfully")

	return nc
}
