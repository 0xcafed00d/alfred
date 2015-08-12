package main

import (
	"os"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	execWithTimeout("go", "env", "/home/lmw/test", os.Stdout, 5*time.Second)
}
