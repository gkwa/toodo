package mdfind

import (
	"fmt"
	"time"
)

type DateCondition struct {
	Param    string
	Duration time.Duration
}

func (d *DateCondition) Slice() []string {
	return []string{d.String()}
}

func (d *DateCondition) String() string {
	return fmt.Sprintf(`%s >= $time.now(-%d)`, d.Param, int64(d.Duration.Seconds()))
}

func NewDateCondition(duration time.Duration) *DateCondition {
	return &DateCondition{
		Param:    "kMDItemFSContentChangeDate",
		Duration: duration,
	}
}
