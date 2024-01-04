package mdfind

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
)

type DirCondition struct {
	Param string
	Dir   string
}

func NewDirCondition(dir string) *DirCondition {
	d := &DirCondition{
		Param: "-onlyin",
		Dir:   dir,
	}

	err := d.Expand()
	if err != nil {
		panic(err)
	}

	return d
}

func (d *DirCondition) Expand() error {
	expanded, err := homedir.Expand(d.Dir)
	if err != nil {
		return err
	}

	d.Dir = expanded
	return nil
}

func (d *DirCondition) Slice() []string {
	if d.Dir == "" {
		return []string{}
	}

	return []string{
		d.Param, d.Dir,
	}
}

func (d *DirCondition) String() string {
	if d.Dir == "" {
		return ""
	}

	return fmt.Sprintf(`%s %s`, d.Param, d.Dir)
}
