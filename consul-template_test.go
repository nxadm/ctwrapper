package main

import "testing"

//import (
//	"os"
//	"testing"
//)
//
func TestCreateTmplArg(t *testing.T) {
	tmplArgConstr := createTmplArg(templatesTest, defaultExt)
	if len(tmplArgConstr) != len(tmplArgs) {
		t.Errorf("Incorrect number of templates, got: %d, want: %d.", len(tmplArgConstr), len(tmplArgs))
	}
	for _, elem := range tmplArgConstr {
		if !stringInSlice(elem, tmplArgs) {
			t.Errorf("\"%s\" not found in the predefined data (found files: %+v)", elem, tmplArgs)
		}
	}
}

// TODO: enable and debug
//func TestRunCt(t *testing.T) {
//	osPath := os.Getenv("PATH")
//	defer os.Setenv("PATH", osPath)
//	os.Setenv("PATH", "t")
//	err := runCt(defaultExt, templatesTest, []string{"-exec", "true"})
//	if err != nil {
//		t.Errorf("Error running the executable: %s", err)
//	}
//}
