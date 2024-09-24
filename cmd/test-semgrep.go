package main

import (
	"context"
	"fmt"
	"os/exec"
)

func main() {
	// Another example of untrusted input
	input := "ping -c 8 google.com; echo hacked"

	ctx := context.Background()

	// Vulnerable: input is directly concatenated into the command
	command := fmt.Sprintf("sh -c %s", input)
	cmd := exec.CommandContext(ctx, command)

	// Execute and print the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Output:", string(output))
}

