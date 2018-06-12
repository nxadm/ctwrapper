package main

import (
	"os"
	"testing"
)

func TestClone(t *testing.T) {
	err := repoTest.clone()
	if err != nil {
		t.Errorf("Unexpected error found: %s", err)
	}

	err = os.RemoveAll(repoTest.Dir)
	if err != nil {
		t.Errorf("Unexpected error found: %s", err)
	}
}
