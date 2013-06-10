package main

import (
	"testing"
)

func TestEncoding(t *testing.T) {
	cmd := "wc -l /bin/date"
	expected := "e2c569be17396eca2a2e3c11578123ed"
	encoded := CreateHash(cmd)

	if encoded != expected {
		t.Fatal("Failed to create hash:", encoded)
	}
}
