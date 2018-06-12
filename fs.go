package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func findConsulTemplate() (string, error) {
	var executable string

	switch {
	// Look in current directory
	case true:
		cwd, err := os.Getwd()
		if err == nil {
			executable = cwd + "/consul-template"
			if _, err := os.Stat(executable); err == nil {
				return executable, nil
			}
		} else {
			return "", err
		}
		fallthrough
	// Look in Path
	default:
		return exec.LookPath("consul-template")

	}

}

func findTemplates(dir, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ext {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}
