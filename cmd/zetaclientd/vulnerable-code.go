package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Print("Enter a command to execute: ")
	var input string
	fmt.Scanln(&input)

	// Vulnerability: Exec without proper validation
	cmd := exec.Command("/bin/sh", "-c", input)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Output:", string(output))
}
