package certmanager

import (
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

// Manager handles certificate generation and management
type Manager struct {
	caCertPath     string
	caKeyPath      string
	serverCertPath string
	serverKeyPath  string
	basePath       string
	caCert         *x509.Certificate
	caKey          *rsa.PrivateKey
}

// NewManager creates a new certificate manager
func NewManager(basePath string) *Manager {
	return &Manager{
		basePath: basePath,
	}
}

// GenerateCACert generates a CA certificate
func (m *Manager) GenerateCACert() error {
	// Generate key pair
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         "AegisGate CA",
			Organization:       []string{"AegisGate"},
			Country:            []string{"US"},
			Province:           []string{"California"},
			Locality:           []string{"San Francisco"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create self-signed certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Parse certificate - FIXED: use cert variable instead of _
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Save files with secure permissions
	m.basePath = "config/certificates"
	if err := os.MkdirAll(m.basePath, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	m.caCertPath = filepath.Join(m.basePath, "ca.crt")
	m.caKeyPath = filepath.Join(m.basePath, "ca.key")

	if err := os.WriteFile(m.caCertPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}), 0600); err != nil {
		return fmt.Errorf("failed to write CA certificate: %w", err)
	}

	// Save private key
	keyDER := x509.MarshalPKCS1PrivateKey(key)
	if err := os.WriteFile(m.caKeyPath, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDER}), 0600); err != nil {
		return fmt.Errorf("failed to write CA key: %w", err)
	}

	// FIXED: assign cert and key to manager
	m.caCert = cert
	m.caKey = key

	return nil
}

// GenerateServerCert generates a server certificate signed by CA
func (m *Manager) GenerateServerCert() error {
	if m.caCert == nil || m.caKey == nil {
		if err := m.GenerateCACert(); err != nil {
			return err
		}
	}

	// Generate key pair
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName:         "localhost",
			Organization:       []string{"AegisGate"},
			Country:            []string{"US"},
		},
		DNSNames:              []string{"localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create signed certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, m.caCert, &key.PublicKey, m.caKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Parse certificate - FIXED: use cert variable instead of _
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Save files with secure permissions
	m.serverCertPath = filepath.Join(m.basePath, "server.crt")
	m.serverKeyPath = filepath.Join(m.basePath, "server.key")

	if err := os.WriteFile(m.serverCertPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}), 0600); err != nil {
		return fmt.Errorf("failed to write server certificate: %w", err)
	}

	// Save private key
	keyDER := x509.MarshalPKCS1PrivateKey(key)
	if err := os.WriteFile(m.serverKeyPath, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDER}), 0600); err != nil {
		return fmt.Errorf("failed to write server key: %w", err)
	}

	// FIXED: assign cert to manager
	m.caCert = cert

	return nil
}

// GetCAPath returns the CA certificate path
func (m *Manager) GetCAPath() string {
	return m.caCertPath
}

// GetServerPath returns the server certificate path
func (m *Manager) GetServerPath() string {
	return m.serverCertPath
}
