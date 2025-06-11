package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const incrementalSearchTool = "peco"

func selectTaskNameCommand(taskfile string) string {
	return fmt.Sprintf(`
	  task -t %s -l --sort none | \
	    tail -n +2 | \
	    sed 's/^\*//' | \
	    %s | \
	    sed -E 's/^ ([^ ]+):.*/\1/' | \
	    sed -E 's/:$//'
	`, taskfile, incrementalSearchTool)
}

func selectTaskName(taskfile string) (string, error) {
	cmd := exec.Command("sh", "-c", selectTaskNameCommand(taskfile))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.Output: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func readFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return file, nil
}

func readOptionalInput(prompt string) string {
	return readInput(prompt, false)
}

func readRequiredInput(prompt string) string {
	return readInput(prompt, true)
}

func readInput(prompt string, required bool) string {
	necessity := "optional"
	if required {
		necessity = "required"
	}

	fmt.Printf(`Enter var "%s" (%s): `, prompt, necessity)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

func runTask(taskfile string, name string, args ...string) error {
	fmt.Printf("run: task -t %s %s %s\n", taskfile, name, strings.Join(args, " "))

	cmdArgs := append([]string{"-t", taskfile, name}, args...)
	cmd := exec.Command("task", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run: %w", err)
	}

	return nil
}
