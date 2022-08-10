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

type Diff interface {
	RenderHumanReadable() string
}

type DependencyDiff struct {
	Type     DiffType
	MetaType DiffMetaType
	Name     string
	Parent   string
	Version  string
}

func (diff DependencyDiff) RenderHumanReadable() string {
	prepos := ""
	if diff.Type == ADDED {
		prepos = "to"
	} else {
		prepos = "of"
	}
	return fmt.Sprintf("%s %s %s %s %s", diff.Type, diff.MetaType, diff.Name, prepos, diff.Parent)
}

type FieldDiff struct {
	Type     DiffType
	MetaType DiffMetaType
	Field    string
	Name     string
	OldValue string
	NewValue string
}

func (diff FieldDiff) RenderHumanReadable() string {
	return fmt.Sprintf("%s %s %s | (old)%s=%s & (new)%s=%s",
		diff.Type,
		diff.MetaType,
		diff.Name,
		diff.Field,
		diff.OldValue,
		diff.Field,
		diff.NewValue)
}

func GenerateDependencyFieldDiff(pkgName string, fieldName string, oldVal string, newVal string) FieldDiff {
	return FieldDiff{
		Type:     MODIFIED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Field:    fieldName,
		OldValue: oldVal,
		NewValue: newVal,
	}
}

func GenerateAddedDependencyDiff(pkgName string, version string) DependencyDiff {
	return DependencyDiff{
		Type:     ADDED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Parent:   "",
		Version:  version,
	}
}

func GenerateRemovedDependencyDiff(pkgName string) DependencyDiff {
	return DependencyDiff{
		Type:     REMOVED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Parent:   "",
	}
}

func GenerateRemovedSubDependencyDiff(pkgName, parent string) DependencyDiff {
	return DependencyDiff{
		Type:     REMOVED,
		MetaType: SUB_DEPENDENCY,
		Name:     pkgName,
		Parent:   parent,
	}
}

func GenerateAddedSubDependencyDiff(pkgName string, parent string, version string) DependencyDiff {
	return DependencyDiff{
		Type:     ADDED,
		MetaType: SUB_DEPENDENCY,
		Name:     pkgName,
		Parent:   parent,
		Version:  version,
	}
}
