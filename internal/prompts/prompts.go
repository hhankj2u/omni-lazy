package prompts

import (
	"log"
	"os/exec"
	"runtime"
	"strings"
)

type Prompt struct {
}

func NewPrompt() *Prompt {
	return &Prompt{}
}

// FetchFabricResult executes the fabric command with the given options and returns the result
func (p *Prompt) FetchFabricResult(command string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
		log.Println("command: ", cmd)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
