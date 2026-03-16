package pci

import (
	"testing"
)

func TestPCIModule_NewModule(t *testing.T) {
	m := NewPCIModule()

	if m == nil {
		t.Fatal("NewPCIModule returned nil")
	}

	if len(m.cardPatterns) == 0 {
		t.Error("No card patterns initialized")
	}
}

func TestPCICardPatterns(t *testing.T) {
	m := NewPCIModule()

	if len(m.cardPatterns) == 0 {
		t.Error("No card patterns initialized")
	}

	t.Logf("Card patterns count: %d", len(m.cardPatterns))
}
