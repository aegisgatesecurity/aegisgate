// +build ignore

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Generate test server key and certificate
	serverKey, serverCert := generateServerCert("aegisgate.local", 1)

	// Write pkg/tls/certs/
	writeKeyAndCert("pkg/tls/certs/server.key", "pkg/tls/certs/server.crt", serverKey, serverCert)

	// Write pkg/adapters/certs/
	writeKeyAndCert("pkg/adapters/certs/server.key", "pkg/adapters/certs/server.crt", serverKey, serverCert)

	// Generate CA key and certificate
	caKey, caCert := generateCA("Test MITM CA", 10)

	// Write pkg/proxy/certs/ca/
	writeKeyAndCert("pkg/proxy/certs/ca/ca.key", "pkg/proxy/certs/ca/ca.crt", caKey, caCert)

	// Write pkg/adapters/custom-certs/
	_, customCert := generateServerCert("custom.example.com", 1)
	writeKeyAndCert("pkg/adapters/custom-certs/server.key", "pkg/adapters/custom-certs/server.crt", serverKey, customCert)

	fmt.Println("All test certificates and keys generated successfully!")
}

func generateServerCert(host string, validYears int) (*rsa.PrivateKey, *x509.Certificate) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate server key: %v", err))
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"AegisGate Test"},
			CommonName:   host,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(validYears, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{host, "localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to create certificate: %v", err))
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse certificate: %v", err))
	}

	return privateKey, cert
}

func generateCA(org string, validYears int) (*ecdsa.PrivateKey, *x509.Certificate) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate CA key: %v", err))
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{org},
			CommonName:   org,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(validYears, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to create CA certificate: %v", err))
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse CA certificate: %v", err))
	}

	return privateKey, cert
}

func writeKeyAndCert(keyPath, certPath string, key interface{}, cert *x509.Certificate) {
	// Write private key
	keyFile, err := os.Create(keyPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to create key file %s: %v", keyPath, err))
	}
	defer keyFile.Close()

	var keyBytes []byte
	switch k := key.(type) {
	case *rsa.PrivateKey:
		keyBytes = x509.MarshalPKCS1PrivateKey(k)
	case *ecdsa.PrivateKey:
		var err error
		keyBytes, err = x509.MarshalECPrivateKey(k)
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal ECDSA key: %v", err))
		}
	}
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}); err != nil {
		panic(fmt.Sprintf("Failed to encode key: %v", err))
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(keyPath), 0755); err != nil {
		panic(fmt.Sprintf("Failed to create directory: %v", err))
	}
	if err := os.MkdirAll(filepath.Dir(certPath), 0755); err != nil {
		panic(fmt.Sprintf("Failed to create directory: %v", err))
	}

	// Write certificate
	certFile, err := os.Create(certPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to create cert file %s: %v", certPath, err))
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
		panic(fmt.Sprintf("Failed to encode cert: %v", err))
	}

	fmt.Printf("Generated: %s and %s\n", keyPath, certPath)
}