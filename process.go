package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func split(s, charset string) []string {
	res := []string{}
	tokenStart := -1

	for i, r := range s {
		if strings.ContainsRune(charset, r) {
			if tokenStart != -1 {
				res = append(res, s[tokenStart:i])
				tokenStart = -1
			}
		} else {
			if tokenStart == -1 {
				tokenStart = i
			}
		}
	}
	if tokenStart != -1 {
		res = append(res, s[tokenStart:])
	}
	return res
}

var ErrTimeout = errors.New("exec timeput")

func execWithTimeout(proc, args, gopath string, out io.Writer, timeout time.Duration) error {

	fmt.Fprintf(out, ">%s %s\n", proc, args)
	defer fmt.Fprintln(out)

	cmd := exec.Command(proc, split(args, " \t")...)
	cmd.Stdout = out
	cmd.Stderr = out

	for _, evar := range os.Environ() {
		if strings.HasPrefix(evar, "GOPATH=") {
			cmd.Env = append(cmd.Env, "GOPATH="+gopath)
		} else {
			cmd.Env = append(cmd.Env, evar)
		}
	}

	c := make(chan error)
	go func(c chan error) {
		c <- cmd.Run()
	}(c)

	timeoutChan := time.NewTimer(timeout)

	select {
	case err := <-c:
		return err
	case <-timeoutChan.C:
		cmd.Process.Kill()
		return ErrTimeout
	}
}
