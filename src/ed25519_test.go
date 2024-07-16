package mpc

import (
	"crypto/ed25519"
	"testing"
)

func TestVanillaEd25519(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Error calling GenerateKey(): %v", err)
	}

	msg := []byte("The quick brown fox jumps over the lazy dog")

	sig := ed25519.Sign(priv, msg)

	if !ed25519.Verify(pub, msg, sig) {
		t.Fatalf("Error verifying signature: %v", err)
	}
}

func TestVanillaEd25519WithOptions(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Error calling GenerateKey(): %v", err)
	}

	msg := []byte("The quick brown fox jumps over the lazy dog")

	sig, err := priv.Sign(nil, msg, &ed25519.Options{
		Context: "Example_ed25519ctx",
	})
	if err != nil {
		t.Fatalf("Error calling Sign(): %v", err)
	}

	if err := ed25519.VerifyWithOptions(pub, msg, sig, &ed25519.Options{
		Context: "Example_ed25519ctx",
	}); err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}
}

func TestVanillaEd25519InvalidateMsg(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Error calling GenerateKey(): %v", err)
	}

	msg := []byte("The quick brown fox jumps over the lazy dog")

	sig, err := priv.Sign(nil, msg, &ed25519.Options{
		Context: "Example_ed25519ctx",
	})
	if err != nil {
		t.Fatalf("Error calling Sign(): %v", err)
	}

	msg = []byte("Big fjords vex quick waltz nymph")

	if err := ed25519.VerifyWithOptions(pub, msg, sig, &ed25519.Options{
		Context: "Example_ed25519ctx",
	}); err == nil {
		t.Fatalf("Expected invalid signature error")
	}
}

func TestVanillaEd25519InvalidateSig(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Error calling GenerateKey(): %v", err)
	}

	msg := []byte("The quick brown fox jumps over the lazy dog")

	sig, err := priv.Sign(nil, msg, &ed25519.Options{
		Context: "Example_ed25519ctx",
	})
	if err != nil {
		t.Fatalf("Error calling Sign(): %v", err)
	}

	sig = append(sig, []byte("expected signature")...)

	if err := ed25519.VerifyWithOptions(pub, msg, sig, &ed25519.Options{
		Context: "Example_ed25519ctx",
	}); err == nil {
		t.Fatalf("Expected invalid signature error")
	}
}
