package main

import (
	"io"
	"os"
	"os/exec"
	"time"
)

func execWithTimeout(proc, args string, env []string, out io.Writer, timeout time.Duration) error {
	cmd, err := exec.Command(proc, args)
	if err != nil {
		return err
	}

	cmd.Stdout = out
	cmd.Stderr = out

}
