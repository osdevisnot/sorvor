// Package cert helps creating a self signed key and certs for TLS server
// Note: Adopted from https://golang.org/src/crypto/tls/generate_cert.go
package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/osdevisnot/sorvor/pkg/logger"
)

// GenerateKeyPair generated self signed key and cert
func GenerateKeyPair(host string) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	logger.Fatal(err, "Failed to create private key")

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	logger.Fatal(err, "Failed to generate serial number")

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"sorvor"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              []string{"localhost", host},
		IPAddresses:           []net.IP{net.IPv6loopback, net.IPv4(127, 0, 0, 1)},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	logger.Fatal(err, "Failed to create certificate")

	certOut, err := os.Create("cert.pem")
	logger.Fatal(err, "Failed to open cert.pem for writing")
	defer certOut.Close()
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	logger.Fatal(err, "Failed to write data to cert.pem")

	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	logger.Fatal(err, "Failed to open key.pem for writing")
	defer keyOut.Close()

	privateBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	logger.Fatal(err, "Unable to marshal private key")

	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privateBytes})
	logger.Fatal(err, "Failed to write data to key.pem")
}
