package toodo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandString(t *testing.T) {
	tests := []struct {
		name     string
		command  *Command
		expected string
	}{
		{
			name: "mdfind command with specific conditions",
			command: &Command{
				Name: "mdfind",
				Arguments: []string{
					`kMDItemFSContentChangeDate >= $time.now(-3600) && (kMDItemFSName == '*.go' || kMDItemFSName == '*.txt')`,
					"-onlyin", "~/pdev",
				},
			},
			expected: `mdfind kMDItemFSContentChangeDate >= $time.now(-3600) && (kMDItemFSName == '*.go' || kMDItemFSName == '*.txt') -onlyin ~/pdev`,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.command.String()
			assert.Equal(t, test.expected, result)
		})
	}
}
