package mdfind

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/taylormonacelli/navyloss"
)

type MDFind struct {
	DirCondition        *DirCondition
	DateCondition       *DateCondition
	ExtensionConditions *FileExtensionsConditions
	Command             exec.Cmd
}

func NewMDFind(dir string, date string, exts []string) *MDFind {
	dur, err := navyloss.DurationFromString(date)
	if err != nil {
		return &MDFind{}
	}

	m := &MDFind{
		DirCondition:        NewDirCondition(dir),
		DateCondition:       NewDateCondition(dur),
		ExtensionConditions: NewFileExtensionsConditions(exts),
	}

	m.SetCommand()
	return m
}

func (m *MDFind) SetCommand() {
	cmd := make([]string, 0)
	cmd = append(cmd, m.DirCondition.Slice()...)

	p := m.DateCondition.String()
	if m.ExtensionConditions.String() != "" {
		p += " && " + m.ExtensionConditions.String()
	}

	z := fmt.Sprintf(`"%s"`, p)
	cmd = append(cmd, z)

	m.Command = *exec.Command("mdfind", cmd...)
}

func (m *MDFind) String() string {
	return m.Command.String()
}

func (m *MDFind) Run() ([]byte, error) {
	// Print the command before running it
	fmt.Println("Running command:", m.Command.String())

	// runs command and captures stdout and stderr and stores
	// them both in separate buffers
	var stdout, stderr bytes.Buffer
	m.Command.Stdout = &stdout
	m.Command.Stderr = &stderr
	err := m.Command.Run()
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s\nstdout: %s", err, stderr.String(), stdout.String())

		slog.Error("Command failed", "command", m.Command.String(), "stderr", stderr.String(), "stdout", stdout.String())

		return nil, errors.New(errorMsg)
	}

	return stdout.Bytes(), nil
}
