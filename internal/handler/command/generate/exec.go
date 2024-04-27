package generate

import (
	"fmt"
	osExec "os/exec"
	"regexp"
)

func exec(command, workingDir string) (string, error) {
	re := regexp.MustCompile(`\s+`)
	tokens := re.Split(command, -1)

	if len(tokens) < 1 {
		return "", fmt.Errorf("invalid command: %v", command)
	}

	commandName := tokens[0]
	args := tokens[1:]

	cmd := osExec.Command(commandName, args...)
	cmd.Dir = workingDir
	output, errCommandExec := cmd.CombinedOutput()
	if errCommandExec != nil {
		return formatOutput(output), errCommandExec
	}

	return formatOutput(output), nil
}

func formatOutput(output []byte) string {
	if len(output) == 0 {
		return ""
	}

	return string(output)
}
