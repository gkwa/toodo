package mdfind

import (
	"fmt"
	"strings"
)

type FileExtensionCondition struct {
	Param string
	Ext   string
}

func (f *FileExtensionCondition) String() string {
	return fmt.Sprintf(`%s == '*.%s'`, f.Param, f.Ext)
}

func NewFileExtensionCondition(ext string) *FileExtensionCondition {
	e := &FileExtensionCondition{
		Param: "kMDItemFSName",
		Ext:   ext,
	}
	return e
}

type FileExtensionsConditions struct {
	Exts []FileExtensionCondition
}

func NewFileExtensionsConditions(exts []string) *FileExtensionsConditions {
	x := make([]FileExtensionCondition, len(exts))
	for i, v := range exts {
		x[i] = *NewFileExtensionCondition(v)
	}

	return &FileExtensionsConditions{
		Exts: x,
	}
}

func (f *FileExtensionsConditions) String() string {
	if len(f.Exts) == 0 {
		return ""
	}

	if len(f.Exts) == 1 {
		return f.Exts[0].String()
	}

	s := make([]string, len(f.Exts))
	for i, v := range f.Exts {
		s[i] = v.String()
	}

	return "(" + strings.Join(s, " || ") + ")"
}
