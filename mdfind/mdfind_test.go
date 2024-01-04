package mdfind

import (
	"testing"
)

func TestMDFindConditionsString(t *testing.T) {
	tests := []struct {
		name             string
		dir              string
		period           string
		fileExtensions   []string
		expectedCommand  string
		expectedErrorMsg string
	}{
		{
			name:            "Happy path",
			dir:             "/path/to/dir",
			period:          "1d",
			fileExtensions:  []string{"txt"},
			expectedCommand: `/usr/bin/mdfind -onlyin /path/to/dir "kMDItemFSContentChangeDate >= $time.now(-86400) && kMDItemFSName == '*.txt'"`,
		},
		{
			name:            "No file extensions",
			dir:             "/path/to/dir",
			period:          "1d",
			fileExtensions:  []string{},
			expectedCommand: `/usr/bin/mdfind -onlyin /path/to/dir "kMDItemFSContentChangeDate >= $time.now(-86400)"`,
		},
		{
			name:            "Multiple file extensions",
			dir:             "/path/to/dir",
			period:          "1d",
			fileExtensions:  []string{"txt", "go"},
			expectedCommand: `/usr/bin/mdfind -onlyin /path/to/dir "kMDItemFSContentChangeDate >= $time.now(-86400) && (kMDItemFSName == '*.txt' || kMDItemFSName == '*.go')"`,
		},
		{
			name:            "Don't pass directory",
			dir:             "",
			period:          "1d",
			fileExtensions:  []string{"txt"},
			expectedCommand: `/usr/bin/mdfind "kMDItemFSContentChangeDate >= $time.now(-86400) && kMDItemFSName == '*.txt'"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdfind := NewMDFind(tt.dir, tt.period, tt.fileExtensions)

			if tt.expectedErrorMsg == "" {
				// Expecting a valid command
				got := mdfind.String()

				if got != tt.expectedCommand {
					t.Errorf("Got: %s, Expected: %s", got, tt.expectedCommand)
				}
			} else {
				// Expecting an error
				err := mdfind.Command.Run()
				if err == nil || err.Error() != tt.expectedErrorMsg {
					t.Errorf("Got error: %v, Expected error: %s", err, tt.expectedErrorMsg)
				}
			}
		})
	}
}
