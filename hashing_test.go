package main

import (
	"testing"
	"fmt"
)

func TestEncoding(t *testing.T) {
	if CreateHash("wc -l /bin/date") != "d2MgLWwgL2Jpbi9kYXRl" {
		t.Fatal("Failed to create hash")
	}

	if DecodeHash("d2MgLWwgL2Jpbi9kYXRl") != "wc -l /bin/date" {
		t.Fatal("Could not decode hash")
	}

	plaintext := "LIKE OMG"
	hashed := CreateHash(plaintext)
	decoded := DecodeHash(hashed)
	if plaintext != decoded {

		fmt.Println([]byte(plaintext))
		fmt.Println([]byte(decoded)[:len(decoded)-1])
		t.Fatal("Encoding was not decodable", "\"" + plaintext + "\"", "vs", "\"" + decoded + "\"")
	}
}
