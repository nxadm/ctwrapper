package main

import "testing"

func TestRetrieveVaultSecret(t *testing.T) {
	_, err := retrieveVaultSecret("https://127.0.0.1", "foo/bar")
	if err == nil {
		t.Error("Expected a error. None thrown.")
	}
}
