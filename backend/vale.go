package backend

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CheckWithVale checks the given text with Vale and returns the output.

func TestVale(text string) (string, error) {
	fmt.Println("TestVale function called with input:", text)
	return "test", nil
}

func CheckWithVale(text string) (string, error) {
	// Create the command.
	cmd := exec.Command("vale", "--no-exit", "--output=CLI", "-")

	// Pass the text to Vale via stdin
	cmd.Stdin = bytes.NewBufferString(text)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s", err)
	}

	return string(output), nil
}
