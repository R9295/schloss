package core

type DiffType int

const (
	ADDED DiffType = iota
	MODIFIED
	REMOVED
)

type DiffMetaType int

const (
	DEPENDENCY DiffMetaType = iota
	SUB_DEPENDENCY
)

type Diff struct {
	// The Type of Diff
	Type     DiffType
	MetaType DiffMetaType
	Name     string
	Text     string
}
