package main

import (
	os/exec
)

func runCommand(cmd string) {
	exec.Command("sh", "-c", cmd).Run()
}
