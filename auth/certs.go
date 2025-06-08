package auth

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"

	"software.sslmate.com/src/go-pkcs12"

	"github.com/BloodRedCodr/go-utils/logger"
)

func GetCertsFromP12(logger *logger.Logger, p12CertPath string, p12Pwd string) (tls.Certificate, *x509.CertPool, string, string) {
	p12Data, err := os.ReadFile(p12CertPath)
	if err != nil {
		logger.Fatal("Failed to read .p12 file: %v", err)
	}

	privateKey, cert, caCerts, err := pkcs12.DecodeChain(p12Data, p12Pwd)
	if err != nil {
		logger.Fatal("Failed to decode .p12 file: %v", err)
	}

	var certPEM, keyPEM []byte

	// PEM encode the certificate chain (leaf + intermediates)
	for _, c := range append([]*x509.Certificate{cert}, caCerts...) {
		block := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: c.Raw,
		}
		certPEM = append(certPEM, pem.EncodeToMemory(block)...)
	}

	// PEM encode the private key
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		logger.Fatal("Failed to marshal private key: %v", err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		logger.Fatal("Failed to create TLS certificate from .p12: %v", err)
	}

	caCertPool := x509.NewCertPool()
	for _, ca := range caCerts {
		caCertPool.AddCert(ca)
	}

	writeTempFile := func(prefix string, data []byte) string {
		tmpFile, err := os.CreateTemp("", prefix)
		if err != nil {
			logger.Fatal("Failed to create temp file: %v", err)
		}
		if _, err := tmpFile.Write(data); err != nil {
			logger.Fatal("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()
		return tmpFile.Name()
	}

	certPath := writeTempFile("cert-", certPEM)
	keyPath := writeTempFile("key-", keyPEM)

	return tlsCert, caCertPool, certPath, keyPath
}
