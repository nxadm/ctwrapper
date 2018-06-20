package main

//func TestImportValues(t *testing.T) {
//	failCli := map[string][]string{
//		"nopass": []string{"ctwrapper", "-r", "foo", "-u", "me"},
//		"norepo": []string{"ctwrapper", "-u", "me", "-p", "secret/foo"},
//	}
//	okCli := map[string][]string{
//		"minimal": []string{"ctwrapper", "-r", "foo", "-d", "/project"},
//	}
//
//	for k, v := range failCli {
//		os.Args = v
//		config := Config{}
//		err := config.importValues()
//		if err == nil {
//			t.Errorf("Expected a error (%s). None thrown", k)
//		}
//	}
//
//	for k, v := range okCli {
//		os.Args = v
//		config := Config{}
//		err := config.importValues()
//		if err != nil {
//			t.Errorf("Error thrown when none expected (%s): %s", k, err)
//		}
//	}
//
//}

//func TestSplitArg(t *testing.T) {
//	givenArgs := []string{
//		`-exec`, `'ls -la'`, `-template`,
//		`"/some/template with spaces:template_no_spaces"`,
//	}
//	arg := strings.Join(givenArgs, " ")
//	args := splitArg(arg)
//	if len(args) != len(givenArgs) {
//		t.Errorf("Incorrect number of args: got %d, expected %d", len(args), len(givenArgs))
//	}
//}
