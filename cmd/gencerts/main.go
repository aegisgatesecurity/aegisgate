package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "math/big"
    "os"
    "time"
)

func main() {
    // Generate CA key
    caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to generate CA key: %v\n", err)
        os.Exit(1)
    }

    // Create CA certificate
    caTemplate := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"AegisGate"},
            CommonName:   "AegisGate CA",
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
        BasicConstraintsValid: true,
        IsCA:                  true,
    }

    caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create CA certificate: %v\n", err)
        os.Exit(1)
    }

    caCert, err := x509.ParseCertificate(caCertDER)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to parse CA certificate: %v\n", err)
        os.Exit(1)
    }

    // Generate client key
    clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to generate client key: %v\n", err)
        os.Exit(1)
    }

    // Create client certificate
    clientTemplate := x509.Certificate{
        SerialNumber: big.NewInt(2),
        Subject: pkix.Name{
            Organization: []string{"AegisGate"},
            CommonName:   "aegisgate-client",
        },
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(90 * 24 * time.Hour),
        KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
        IPAddresses:  nil,
    }

    clientCertDER, err := x509.CreateCertificate(rand.Reader, &clientTemplate, caCert, &clientKey.PublicKey, caKey)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create client certificate: %v\n", err)
        os.Exit(1)
    }

    // Write CA certificate
    caCertFile, err := os.Create("ca.crt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create ca.crt: %v\n", err)
        os.Exit(1)
    }
    defer func() {
        if err := caCertFile.Close(); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to close ca.crt: %v\n", err)
        }
    }()
    if err := pem.Encode(caCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertDER}); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode CA certificate: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Generated: ca.crt")

    // Write client certificate
    clientCertFile, err := os.Create("client.crt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create client.crt: %v\n", err)
        os.Exit(1)
    }
    defer func() {
        if err := clientCertFile.Close(); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to close client.crt: %v\n", err)
        }
    }()
    if err := pem.Encode(clientCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: clientCertDER}); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode client certificate: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Generated: client.crt")

    // Write client key
    clientKeyFile, err := os.Create("client.key")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create client.key: %v\n", err)
        os.Exit(1)
    }
    defer func() {
        if err := clientKeyFile.Close(); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to close client.key: %v\n", err)
        }
    }()
    clientKeyBytes, err := x509.MarshalECPrivateKey(clientKey)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to marshal EC private key: %v\n", err)
        os.Exit(1)
    }
    if err := pem.Encode(clientKeyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: clientKeyBytes}); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode client key: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Generated: client.key")

    fmt.Println("\nmTLS certificates generated successfully!")
    fmt.Println("Files: ca.crt, client.crt, client.key")
}
