package crypto

import "testing"

func init() {
	// Lower the work factor so the suite runs fast.
	SetWorkFactor(10)
}

func TestRoundTrip(t *testing.T) {
	plain := []byte(`{"version":1,"projects":[]}`)
	ct, err := Encrypt(plain, "hunter2")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	pt, err := Decrypt(ct, "hunter2")
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(pt) != string(plain) {
		t.Fatalf("mismatch: got %q want %q", pt, plain)
	}
}

func TestWrongPassphrase(t *testing.T) {
	ct, err := Encrypt([]byte("secret"), "right")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if _, err := Decrypt(ct, "wrong"); err == nil {
		t.Fatal("expected decrypt to fail with wrong passphrase")
	}
}
