package core

import "fmt"

type DiffType string

const (
	ADDED    DiffType = "added"
	MODIFIED          = "modified"
	REMOVED           = "removed"
)

type DiffMetaType string

const (
	DEPENDENCY     DiffMetaType = "dependency"
	SUB_DEPENDENCY              = "sub-dependency"
)

type Diff struct {
	// The Type of Diff
	Type     DiffType
	MetaType DiffMetaType
	Name     string
	Text     string
}

func GenerateDependencyFieldDiff(pkgName string, fieldName string, oldVal string, newVal string) Diff {
	return Diff{
		Type:     MODIFIED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Text:     fmt.Sprintf("(old)%s=%s & (new)%s=%s", fieldName, oldVal, fieldName, newVal),
	}
}

func GenerateAddedDependencyDiff(pkgName string, version string) Diff {
	return Diff{
		Type:     ADDED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Text:     fmt.Sprintf("version=%s", version),
	}
}

func GenerateRemovedDependencyDiff(pkgName string) Diff {
	return Diff{
		Type:     REMOVED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
	}
}

func GenerateRemovedSubDependencyDiff(pkgName, of string) Diff {
	return Diff{
		Type:     REMOVED,
		MetaType: SUB_DEPENDENCY,
		Name:     pkgName,
		Text:     fmt.Sprintf("of %s", of),
	}
}

func GenerateAddedSubDependencyDiff(pkgName string, to string, version string) Diff {
	return Diff{
		Type:     ADDED,
		MetaType: SUB_DEPENDENCY,
		Name:     pkgName,
		Text:     fmt.Sprintf("to %s", to),
	}
}
