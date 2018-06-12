package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindConsulTemplate(t *testing.T) {
	// Manipulate the PATH
	osPath := os.Getenv("PATH")
	cwd, _ := os.Getwd()
	defer os.Setenv("PATH", osPath)
	os.Setenv("PATH", "")

	// Test exec in cwd
	os.Chdir("t")
	ct, err := findConsulTemplate()
	if err != nil {
		t.Error("Can't find executable in the working directory")
	} else if filepath.Base(ct) != "consul-template" {
		t.Errorf("Wrong executable found in the working directory; %s", ct)
	}

	// Test exec in PATH
	os.Chdir(cwd)
	os.Setenv("PATH", "t")
	if err != nil {
		t.Error("Can't find executable in the PATH")
	} else if filepath.Base(ct) != "consul-template" {
		t.Errorf("Wrong executable found in the working directory; %s", ct)
	}

}

func TestFindTemplates(t *testing.T) {
	templates, err := findTemplates("t", defaultExt)
	if err != nil {
		t.Errorf("Unexpected error found: %s", err)
	}

	if len(templates) != len(templatesTest) {
		t.Errorf("Incorrect number of templates, got: %d, want: %d.", len(templates), len(templatesTest))
	}
	for _, elem := range templatesTest {
		if !stringInSlice(elem, templates) {
			t.Errorf("\"%s\" not found in the predefined data (found files: %+v)", elem, templates)
		}
	}
}
