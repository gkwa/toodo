package mdfind

// create interface that contains slice ans string methods
type Condition interface {
	Slice() []string
	String() string
}
