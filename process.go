package main

import (
	"io"
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

func execWithTimeout(proc, args string, env []string, out io.Writer, timeout time.Duration) error {

	cmd := exec.Command(proc, split(args, " \t")...)
	cmd.Stdout = out
	cmd.Stderr = out

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
		return nil
	}
}
