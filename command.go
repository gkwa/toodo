package toodo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
)

type Command struct {
	Name      string   // Command name (e.g., "mdfind")
	Arguments []string // Command arguments
	Stdout    string   // Captured stdout
	Stderr    string   // Captured stderr
}

func NewCommand(name string) *Command {
	return &Command{
		Name: name,
	}
}

func (c *Command) AddArgument(arg ...string) {
	c.Arguments = append(c.Arguments, arg...)
}

func (c *Command) Run() (string, error) {
	c.ExpandHomeDir()
	cmd := exec.Command(c.Name, c.Arguments...)
	cmd.Env = os.Environ()

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		c.Stderr = strings.TrimSpace(string(stdoutStderr))

		if c.Stderr != "" {
			fmt.Printf("Command stderr: %s\n", c.Stderr)
		}

		return "", fmt.Errorf("command failed: %v", c.Stderr)
	}

	c.Stdout = strings.TrimSpace(string(stdoutStderr))
	return c.Stdout, nil
}

func (c *Command) GetStdout() string {
	return c.Stdout
}

func (c *Command) GetStderr() string {
	return c.Stderr
}

func (c *Command) ExpandHomeDir() {
	for i, arg := range c.Arguments {
		expanded, err := homedir.Expand(arg)
		if err == nil {
			c.Arguments[i] = expanded
		}
	}
}

func (c *Command) String() string {
	tmpl, err := template.New("command").Parse(c.Name + " {{range .Arguments}}{{.}} {{end}}")
	if err != nil {
		return err.Error()
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, struct{ Arguments []string }{c.Arguments}); err != nil {
		return err.Error()
	}

	return strings.TrimSpace(string(output.String()))
}
