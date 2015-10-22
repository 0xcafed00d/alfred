package main

import (
	"os"
)

func runCommand(cmd string, args []string) {
	switch cmd {
	case "serve":
		cmdServe(args)
	case "kick":
		cmdKick(args)
	}
}

func main() {
	if len(os.Args) > 1 {
		runCommand(os.Args[1], os.Args[2:])
	}
}
