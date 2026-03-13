package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Tier int
const (
	TierCommunity    Tier = 0
	TierDeveloper    Tier = 1
	TierProfessional Tier = 2
	TierEnterprise   Tier = 3
)

type LicenseType string
const (
	LicenseTypeCommunity    LicenseType = "community"
	LicenseTypeDeveloper    LicenseType = "developer"
	LicenseTypeProfessional LicenseType = "professional"
	LicenseTypeEnterprise   LicenseType = "enterprise"
)

type License struct {
	ID           string      `json:"id"`
	Type         LicenseType `json:"type"`
	Email        string      `json:"email"`
	Organization string      `json:"organization,omitempty"`
	Tiers        []Tier      `json:"tiers,omitempty"`
	IssuedAt     time.Time   `json:"issued_at"`
	ExpiresAt    time.Time   `json:"expires_at"`
	HardwareID   string      `json:"hardware_id,omitempty"`
	Signature    string      `json:"signature"`
}

func main() {
	var (
		tier       = flag.String("tier", "community", "License tier: community, developer, professional, enterprise")
		email      = flag.String("email", "", "License holder email (required)")
		org        = flag.String("org", "", "Organization name")
		days       = flag.Int("days", 365, "License validity in days")
		hardwareID = flag.String("hardware", "", "Hardware ID (Enterprise only)")
		secretPath = flag.String("secret", "", "Path to HMAC secret file (required)")
	)
	flag.Parse()

	if *email == "" || *secretPath == "" {
		fmt.Fprintf(os.Stderr, "Usage: licensegen -email=user@example.com -secret=/path/to/secret [-tier=professional] [-days=365]\n")
		fmt.Fprintf(os.Stderr, "\nRequired flags:\n")
		fmt.Fprintf(os.Stderr, "  -email    License holder email\n")
		fmt.Fprintf(os.Stderr, "  -secret   Path to HMAC secret file\n")
		fmt.Fprintf(os.Stderr, "\nOptional flags:\n")
		fmt.Fprintf(os.Stderr, "  -tier     License tier (default: community)\n")
		fmt.Fprintf(os.Stderr, "  -days     Validity in days (default: 365)\n")
		fmt.Fprintf(os.Stderr, "  -org      Organization name\n")
		fmt.Fprintf(os.Stderr, "  -hardware Hardware ID (Enterprise tier only)\n")
		os.Exit(1)
	}

	// Load HMAC secret
	secret, err := os.ReadFile(*secretPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading secret file: %v\n", err)
		os.Exit(1)
	}

	// Determine license type and tiers
	var licType LicenseType
	var tiers []Tier
	switch *tier {
	case "community":
		licType = LicenseTypeCommunity
		tiers = []Tier{TierCommunity}
	case "developer":
		licType = LicenseTypeDeveloper
		tiers = []Tier{TierCommunity, TierDeveloper}
	case "professional":
		licType = LicenseTypeProfessional
		tiers = []Tier{TierCommunity, TierDeveloper, TierProfessional}
	case "enterprise":
		licType = LicenseTypeEnterprise
		tiers = []Tier{TierCommunity, TierDeveloper, TierProfessional, TierEnterprise}
	default:
		fmt.Fprintf(os.Stderr, "Invalid tier: %s\n", *tier)
		fmt.Fprintf(os.Stderr, "Valid tiers: community, developer, professional, enterprise\n")
		os.Exit(1)
	}

	// Generate license
	license := License{
		ID:           generateLicenseID(),
		Type:         licType,
		Email:        *email,
		Organization: *org,
		Tiers:        tiers,
		IssuedAt:     time.Now(),
		ExpiresAt:    time.Now().Add(time.Duration(*days) * 24 * time.Hour),
		HardwareID:   *hardwareID,
	}

	// Sign license
	payload, _ := json.Marshal(license)
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	// Encode payload
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	licenseKey := fmt.Sprintf("%s.%s", encodedPayload, signature)

	// Output
	fmt.Println("========================================")
	fmt.Println("     AEGISGATE LICENSE GENERATED")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("License Key:")
	fmt.Println(licenseKey)
	fmt.Println()
	fmt.Println("----------------------------------------")
	fmt.Printf("Tier:        %s\n", *tier)
	fmt.Printf("Email:       %s\n", *email)
	if *org != "" {
		fmt.Printf("Organization: %s\n", *org)
	}
	fmt.Printf("Expires:     %s\n", license.ExpiresAt.Format("2006-01-02"))
	fmt.Printf("License ID:  %s\n", license.ID)
	if *hardwareID != "" {
		fmt.Printf("Hardware ID: %s\n", *hardwareID)
		fmt.Println()
		fmt.Println("NOTE: This license is bound to the specific hardware.")
		fmt.Println("      It will NOT work on other machines.")
	} else {
		fmt.Println()
		fmt.Println("NOTE: This license can be used on any compatible hardware.")
	}
	fmt.Println("----------------------------------------")
	fmt.Println()
	fmt.Println("To use this license:")
	fmt.Println("1. Set as environment variable: AEGISGATE_LICENSE_KEY=<key>")
	fmt.Println("2. Or add to your config file")
	fmt.Println("3. Ensure HMAC secret is available for validation")
}

func generateLicenseID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
