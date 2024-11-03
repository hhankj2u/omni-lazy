package prompts

import (
	"os/exec"
	"strings"
)

type Prompt struct {
}

func NewPrompt() *Prompt {
	return &Prompt{}
}

// FetchFabricResult executes the fabric command with the given options and returns the result
func (p *Prompt) FetchFabricResult(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// ListPatterns fetches the list of patterns from fabric
func (p *Prompt) ListPatterns() ([]string, error) {
	cmd := exec.Command("fabric", "--listpatterns")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	patterns := strings.Split(strings.TrimSpace(string(output)), "\n")
	return patterns, nil
}
