package main

import (
	"testing"
)

func TestEncoding(t *testing.T) {
	cmd := "wc -l /bin/date"
	expected := "4841a561915ae3d73e4b4ffb6fcbe630"
	encoded := CreateHash(cmd)

	if encoded != expected {
		t.Fatal("Failed to create hash:", encoded)
	}
}
