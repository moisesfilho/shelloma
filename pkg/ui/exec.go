package ui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"shelloma/pkg/i18n"
)

func ExecuteCommand(cmdStr string, t i18n.Translations) (int, string, error) {
	fmt.Printf("%s%s⚡ %s%s\n\n", Bold, Green, t.Executing, Reset)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		shell := os.Getenv("COMSPEC")
		if shell != "" && strings.Contains(strings.ToLower(shell), "cmd.exe") {
			cmd = exec.Command("cmd.exe", "/C", cmdStr)
		} else {
			cmd = exec.Command("powershell.exe", "-NoProfile", "-Command", cmdStr)
		}
	} else {
		shell := os.Getenv("SHELL")
		if shell == "" {
			if runtime.GOOS == "darwin" {
				shell = "/bin/zsh"
			} else {
				shell = "/bin/sh"
			}
		}
		cmd = exec.Command(shell, "-c", cmdStr)
	}

	var outBuf bytes.Buffer
	multiOut := io.MultiWriter(os.Stdout, &outBuf)
	multiErr := io.MultiWriter(os.Stderr, &outBuf)

	cmd.Stdin = os.Stdin
	cmd.Stdout = multiOut
	cmd.Stderr = multiErr

	err := cmd.Run()
	fmt.Println()

	outputStr := strings.TrimSpace(outBuf.String())
	exitCode := 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return exitCode, outputStr, err
}
