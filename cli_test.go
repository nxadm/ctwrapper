package main

import (
	"os"
	"testing"
)

func TestImportValues(t *testing.T) {
	failCli := map[string][]string{
		"nopass": []string{"ctwrapper", "-r", "foo", "-u", "me"},
		"nopath": []string{"ctwrapper", "-r", "foo", "-u", "me", "-vk", "secret"},
	}
	for k, v := range failCli {
		os.Args = v
		config := Config{}
		err := config.importValues()
		if err == nil {
			t.Errorf("Expected a error (%s). None thrown", k)
		}
	}
}
