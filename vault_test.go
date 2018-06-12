package main

import "testing"

func TestRetrieveVaultSecret(t *testing.T) {
	_, err := retrieveVaultSecret("foo", "bar")
	if err == nil {
		t.Error("Expected a error. None thrown.")
	}
}
