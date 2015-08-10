package main

import (
	"os"
	"testing"
)

func TestExec(t *testing.T) {
	execWithTimeout("dummyproc", "10 a s d f g", nil, os.Stdout, 100)
}
