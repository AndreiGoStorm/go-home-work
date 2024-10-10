package main

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

const (
	success = 0
	failure = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if returnCode = addEnv(env); returnCode != 0 {
		return
	}
	name, args := cmd[0], cmd[1:]
	command := exec.Command(name, args...)
	if returnCode = run(command); returnCode != 0 {
		return
	}
	return
}

func addEnv(env Environment) int {
	for key, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return int(syscall.EINVAL)
			}
			continue
		}
		err := os.Setenv(key, value.Value)
		if err != nil {
			return int(syscall.EINVAL)
		}
	}

	return success
}

func run(command *exec.Cmd) int {
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return failure
	}
	return success
}
