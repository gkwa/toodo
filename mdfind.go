package toodo

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

type MDFind struct {
	OnlyIn     string
	Conditions []string
}

func NewMDFind(dir string, period time.Duration, fileExtensions []string) *MDFind {
	conditions := []string{
		fmt.Sprintf(`kMDItemFSContentChangeDate >= $time.now(-%d)`, int64(period.Seconds())),
	}

	if len(fileExtensions) == 1 {
		// Single extension, use &&
		conditions = append(conditions, fmt.Sprintf(`kMDItemFSName == '*.%s'`, fileExtensions[0]))
	} else if len(fileExtensions) > 1 {
		// Multiple extensions, use ||
		var extConditions []string
		for _, ext := range fileExtensions {
			extConditions = append(extConditions, fmt.Sprintf(`kMDItemFSName == '*.%s'`, ext))
		}
		conditions = append(conditions, "("+strings.Join(extConditions, " || ")+")")
	}

	return &MDFind{
		OnlyIn:     dir,
		Conditions: conditions,
	}
}

func (m *MDFind) ExpandHomeDir() {
	expanded, err := homedir.Expand(m.OnlyIn)
	if err == nil {
		m.OnlyIn = expanded
	}
}

func (m *MDFind) BuildCommand() *Command {
	cmd := NewCommand("mdfind")

	quotedConditions := strings.Join(m.Conditions, " && ")
	cmd.AddArgument(quotedConditions)

	return cmd
}
