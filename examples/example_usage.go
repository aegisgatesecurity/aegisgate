// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================
//
// This file contains proprietary trade secret information.
// Unauthorized reproduction, distribution, or reverse engineering is prohibited.
// =========================================================================

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/pkiattest"
)

func main() {
	slog.Info("PKI Attestation Framework Example")

	// Step 1: Create PKI attestation instance
	attestation, _ := pkiattest.NewAttestation(nil)
	slog.Info("Created attestation instance")

	// Step 2: Generate test CA certificate
	caCert, caKey := generateTestCA()
	slog.Info("Generated test CA certificate")

	// Step 3: Add CA as trust anchor
	anchor, err := attestation.AddTrustAnchor(caCert)
	if err != nil {
		slog.Error("Failed to add trust anchor", "error", err)
		os.Exit(1)
	}
	slog.Info("Added trust anchor", "id", anchor.CertificateID, "subject", caCert.Subject.CommonName)

	// Step 4: Generate leaf certificate signed by CA
	leafCert, leafKey := generateTestLeafCertificate(caCert, caKey)
	slog.Info("Generated leaf certificate", "subject", leafCert.Subject.CommonName)

	// Step 5: Attest the leaf certificate
	result, err := attestation.AttestCertificate(leafCert)
	if err != nil {
		slog.Error("Attestation failed", "error", err)
		os.Exit(1)
	}

	if result.Valid {
		slog.Info("Certificate attestation successful", "reason", result.Reason)
	} else {
		slog.Error("Certificate attestation failed", "reason", result.Reason)
		os.Exit(1)
	}

	// Step 6: Create trust store
	trustStore := pkiattest.NewTrustStore()
	_, err = trustStore.AddTrustAnchor(caCert)
	if err != nil {
		slog.Error("Failed to add trust anchor to store", "error", err)
		os.Exit(1)
	}
	slog.Info("Created trust store")

	// Step 7: Backdoor prevention - use standalone function, not Attestation field
	backdoorPrevention := pkiattest.NewBackdoorPrevention()
	backdoorPrevention.SetAttestation(attestation)
	isBackdoor, reason, err := backdoorPrevention.DetectBackdoor(leafCert)
	if err != nil {
		slog.Error("Backdoor detection failed", "error", err)
		os.Exit(1)
	}

	if isBackdoor {
		slog.Warn("Backdoor detected", "reason", reason)
	} else {
		slog.Info("Backdoor prevention passed", "reason", reason)
	}

	// Step 8: Verify certificate chain using standalone function
	chain, err := pkiattest.VerifyCertificateChain(leafCert, []*pkiattest.TrustAnchor{anchor})
	if err != nil {
		slog.Error("Chain verification failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Certificate chain verified", "chain_length", len(chain))

	// Step 9: Get certificate fingerprint
	fingerprint := pkiattest.CertificateFingerprint(leafCert)
	slog.Info("Certificate fingerprint", "fingerprint", fingerprint)

	// Step 10: PEM encode/decode round-trip
	pemData, err := encodeCertificateToPEM(leafCert)
	if err != nil {
		slog.Error("PEM encode failed", "error", err)
		os.Exit(1)
	}
	decodedCert, err := decodePEMToCertificate(pemData)
	if err != nil {
		slog.Error("PEM decode failed", "error", err)
		os.Exit(1)
	}
	slog.Info("PEM encode/decode successful", "subject", decodedCert.Subject.CommonName)

	// Step 11: Simulate a backdoor certificate (zero serial)
	zeroSerialCert := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName: "backdoor.example.com",
		},
		NotBefore:   time.Now().Add(-24 * time.Hour),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour),
		IsCA:        false,
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	isBackdoor, reason, err = backdoorPrevention.DetectBackdoor(zeroSerialCert)
	if err != nil {
		slog.Error("Backdoor detection failed", "error", err)
		os.Exit(1)
	}

	if isBackdoor {
		slog.Info("Zero serial backdoor detected correctly", "reason", reason)
	}

	// Step 12: Certificate revocation example
	_ = trustStore.AddRevokedCertificate("99999", "revoked.example.com", "compromised key")
	isRevoked, reason, err := trustStore.IsRevoked("99999")
	if err != nil {
		slog.Error("Revocation check failed", "error", err)
		os.Exit(1)
	}

	if isRevoked {
		slog.Info("Revoked certificate detected", "reason", reason)
	}

	// Step 13: Save certificates to file
	saveCertificates(caCert, leafCert, caKey, leafKey)
	slog.Info("Certificates saved to files")

	slog.Info("PKI Attestation Framework Example completed successfully")
}

func encodeCertificateToPEM(cert *x509.Certificate) ([]byte, error) {
	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(certPEM), nil
}

func decodePEMToCertificate(pemData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM data")
	}
	return x509.ParseCertificate(block.Bytes)
}

func generateTestCA() (*x509.Certificate, *rsa.PrivateKey) {
	// Generate CA key
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate CA key: %v", err))
	}

	// Create CA certificate
	caCert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"Example Org"},
			CommonName:   "Example CA",
		},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		// MaxIssuerPathLen removed - not a valid field in x509.Certificate
	}

	// Sign CA certificate
	caDER, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caKey.PublicKey, caKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to create CA certificate: %v", err))
	}

	caCert, err = x509.ParseCertificate(caDER)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse CA certificate: %v", err))
	}

	return caCert, caKey
}

func generateTestLeafCertificate(caCert *x509.Certificate, caKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey) {
	// Generate leaf key
	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate leaf key: %v", err))
	}

	// Create leaf certificate
	leafCert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"Example Org"},
			CommonName:   "example.com",
		},
		DNSNames:              []string{"example.com", "www.example.com"},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Sign leaf certificate with CA
	leafDER, err := x509.CreateCertificate(rand.Reader, leafCert, caCert, &leafKey.PublicKey, caKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to create leaf certificate: %v", err))
	}

	leafCert, err = x509.ParseCertificate(leafDER)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse leaf certificate: %v", err))
	}

	return leafCert, leafKey
}

func saveCertificates(caCert, leafCert *x509.Certificate, caKey, leafKey *rsa.PrivateKey) {
	// Save CA certificate
	caCertPEM, err := encodeCertificateToPEM(caCert)
	if err != nil {
		slog.Error("Failed to encode CA certificate", "error", err)
		return
	}
	err = os.WriteFile("example_ca.pem", caCertPEM, 0644)
	if err != nil {
		slog.Error("Failed to write CA certificate", "error", err)
	}

	// Save leaf certificate
	leafCertPEM, err := encodeCertificateToPEM(leafCert)
	if err != nil {
		slog.Error("Failed to encode leaf certificate", "error", err)
		return
	}
	err = os.WriteFile("example_leaf.pem", leafCertPEM, 0644)
	if err != nil {
		slog.Error("Failed to write leaf certificate", "error", err)
	}

	// Save CA key
	caKeyBytes := x509.MarshalPKCS1PrivateKey(caKey)
	caKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: caKeyBytes,
	}
	caKeyData := pem.EncodeToMemory(caKeyPEM)
	err = os.WriteFile("example_ca.key", caKeyData, 0600)
	if err != nil {
		slog.Error("Failed to write CA key", "error", err)
	}

	// Save leaf key
	leafKeyBytes := x509.MarshalPKCS1PrivateKey(leafKey)
	leafKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: leafKeyBytes,
	}
	leafKeyData := pem.EncodeToMemory(leafKeyPEM)
	err = os.WriteFile("example_leaf.key", leafKeyData, 0600)
	if err != nil {
		slog.Error("Failed to write leaf key", "error", err)
	}

	slog.Info("Certificates saved", "ca_cert", "example_ca.pem", "leaf_cert", "example_leaf.pem", "ca_key", "example_ca.key", "leaf_key", "example_leaf.key")
}
