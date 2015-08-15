package main

import (
	"github.com/simulatedsimian/assert"
	"os"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	assert.Nil(t, execWithTimeout("dummyproc", "", "", os.Stdout, 100*time.Second))
	assert.Nil(t, execWithTimeout("dummyproc", "5", "", os.Stdout, 100*time.Second))
	assert.NotNil(t, execWithTimeout("dummyproc", "5", "", os.Stdout, 1*time.Second))
	assert.NotNil(t, execWithTimeout("dummyproc", "5 1", "", os.Stdout, 100*time.Second))
}
