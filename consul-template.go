package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

func runCt(tmplExt string, files, options []string) error {
	ct, err := findConsulTemplate()
	if err != nil {
		return err
	}

	args := createTmplArg(files, tmplExt)
	args = append(args, options...)
	fmt.Printf("Executing: %s %s\n", ct, strings.Join(args, " "))

	/* Execute the command */
	switch {
	// Windows does not really forks the process: we need a wrapper to write the stdOut end stdErr at the same time
	case runtime.GOOS == "windows":
		return runWhilePrinting(exec.Command(ct, args...))
	default:
		return syscall.Exec(ct, append([]string{ct}, args...), os.Environ())
	}

	return nil
}

func createTmplArg(files []string, tmplExt string) []string {
	tmplArgs := []string{}
	for _, tmpl := range files {
		file := strings.TrimSuffix(tmpl, tmplExt)
		tmplArgs = append(tmplArgs, "-template")
		tmplArgs = append(tmplArgs, tmpl+":"+file)
	}
	return tmplArgs
}

func runWhilePrinting(cmd *exec.Cmd) error {
	var err error

	/* Get stdout and stderr realtime */
	var wg sync.WaitGroup
	wg.Add(2) // StdOut + StdErr
	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	} else {
		scanner := bufio.NewScanner(cmdOutReader)
		go func() {
			defer wg.Done()
			for scanner.Scan() {
				fmt.Fprintf(os.Stdout, "%s\n", scanner.Text())
			}
		}()
	}
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	} else {
		scanner := bufio.NewScanner(cmdErrReader)
		go func() {
			defer wg.Done()
			for scanner.Scan() {
				fmt.Fprintf(os.Stderr, "%s\n", scanner.Text())
			}
		}()
	}

	/* Run the cmd */
	err = cmd.Start()
	if err != nil {
		return err
	}

	waitErr := cmd.Wait()
	if waitErr != nil {
		return waitErr
	}

	wg.Wait() // Wait for StdOut + StdErr
	return err
}
