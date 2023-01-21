package util

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func ExecuteCommand(command string, writers ...io.Writer) error {
	cmd := exec.Command("sh", "-c", command)

	writers = append(writers, os.Stdout, os.Stderr)

	cmd.Stdout = io.MultiWriter(writers...)
	cmd.Stderr = io.MultiWriter(writers...)

	err := cmd.Run()

	return err
}

func ExecuteCommandWithOutput(command string) (string, error) {
	// create io.writer to the string
	var out bytes.Buffer
	err := ExecuteCommand(command, &out)

	if err != nil {
		return "", err
	}

	return out.String(), nil
}
