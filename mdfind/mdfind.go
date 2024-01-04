package mdfind

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/taylormonacelli/navyloss"
)

type MDFind struct {
	DirCondition        *DirCondition
	DateCondition       *DateCondition
	ExtensionConditions *FileExtensionsConditions
	CommandFile         string
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

	cmd = append(cmd, "mdfind")
	cmd = append(cmd, m.DirCondition.Slice()...)

	p := m.DateCondition.String()
	if m.ExtensionConditions.String() != "" {
		p += " && " + m.ExtensionConditions.String()
	}

	z := fmt.Sprintf(`"%s"`, p)
	cmd = append(cmd, z)

	m.Command = *exec.Command(cmd[0], cmd[1:]...)
}

func (m *MDFind) WriteCommandToFile() string {
	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("mdfind_command_%d.sh", time.Now().Unix())
	filePath := filepath.Join(tmpDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("#!/usr/bin/env bash\n")
	if err != nil {
		panic(err)
	}

	str := m.String() + "\n"

	// Replace $ with \$
	str = strings.Replace(str, "$", "\\$", -1)

	_, err = file.WriteString(str)
	if err != nil {
		panic(err)
	}

	err = os.Chmod(filePath, 0o755) // Set execute bit
	if err != nil {
		panic(err)
	}

	return filePath
}

func (m *MDFind) String() string {
	return m.Command.String()
}

func (m *MDFind) Run() ([]byte, error) {
	path := m.WriteCommandToFile()

	slog.Debug("Running script", "path", path)

	cmd := exec.Command("bash", "-c", path)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}

	if stderr.String() == "" {
		fmt.Print(stdout.String())
		os.Remove(path)
	} else {
		fmt.Println("stderr:", stderr.String())
	}

	return stdout.Bytes(), err
}
