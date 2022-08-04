package core

type DiffType string

const (
	ADDED DiffType = "added"
	MODIFIED = "modified"
	REMOVED = "removed"
)

type DiffMetaType string

const (
	DEPENDENCY DiffMetaType = "dependency"
	SUB_DEPENDENCY = "sub-dependency"
)

type Diff struct {
	// The Type of Diff
	Type     DiffType
	MetaType DiffMetaType
	Name     string
	Text     string
}
