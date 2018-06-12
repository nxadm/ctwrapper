package main

import "os"

var repoTest = Repo{
	URL:    "https://github.com/nxadm/ldifdiff.git",
	Branch: "master",
	Commit: "e004ca26f763892a445e7082a9772eea1a1ae673",
	Dir:    os.TempDir() + "/ctwrapper_test",
}
var templatesTest = []string{"t/1.tmp.tmpl", "t/2/2.tmp.tmpl", "t/2/3/3.tmp.tmpl"}
var tmplArgs = []string{
	"-template", "t/1.tmp.tmpl:t/1.tmp",
	"-template", "t/2/2.tmp.tmpl:t/2/2.tmp",
	"-template", "t/2/3/3.tmp.tmpl:t/2/3/3.tmp",
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
