package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	// Simulate the exact encoding
	jsonStr := `{"id":"PAD-TEST","type":"enterprise","email":"test@example.com","modules":["*"],"tiers":["enterprise"],"issued_at":"2026-02-21T18:29:48.7419274-06:00","expires_at":"2027-02-21T18:29:48.7419274-06:00"}`

	encoded := base64.StdEncoding.EncodeToString([]byte(jsonStr))
	fmt.Printf("Encoded: %s\n\n", encoded[:80])

	// Now decode exactly as license.go does
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return
	}

	fmt.Printf("Decoded length: %d\n", len(decoded))
	fmt.Printf("Decoded string: %s\n\n", string(decoded))

	// Check SplitN behavior
	parts := strings.SplitN(string(decoded), ".", 2)
	fmt.Printf("Parts count: %d\n", len(parts))

	var jsonPayload []byte
	if len(parts) == 2 {
		fmt.Println("Using parts[0] as payload (signed license)")
		jsonPayload = []byte(parts[0])
	} else {
		fmt.Println("Using decoded as payload (unsigned license)")
		jsonPayload = decoded
	}

	fmt.Printf("JSON payload length: %d\n", len(jsonPayload))
	fmt.Printf("JSON payload: %s\n\n", string(jsonPayload))

	// Now try to unmarshal
	var result map[string]interface{}
	if err := json.Unmarshal(jsonPayload, &result); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
	} else {
		fmt.Printf("Unmarshal success: %v\n", result)
	}
}
