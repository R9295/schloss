package core

import "fmt"

type DiffType string

const (
	ADDED    DiffType = "added"
	MODIFIED DiffType = "modified"
	REMOVED  DiffType = "removed"
)

type DiffMetaType string

const (
	DEPENDENCY     DiffMetaType = "dependency"
	SUB_DEPENDENCY DiffMetaType = "sub-dependency"
)

type Diff interface {
	RenderHumanReadable() string
	GetType() DiffType
	GetName() string
}

type DependencyDiff struct {
	Type     DiffType     `json:"type"`
	MetaType DiffMetaType `json:"meta_type"`
	Name     string       `json:"dependency_name"`
	Parent   string       `json:"parent"`
	Version  string       `json:"version,omitempty"`
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

func (diff DependencyDiff) GetType() DiffType {
	return diff.Type
}

func (diff DependencyDiff) GetName() string {
	return diff.Name
}

type FieldDiff struct {
	Type     DiffType     `json:"type"`
	MetaType DiffMetaType `json:"meta_type"`
	Name     string       `json:"dependency_name"`
	Field    string       `json:"field"`
	OldValue string       `json:"old_value"`
	NewValue string       `json:"new_value"`
}

func (diff FieldDiff) RenderHumanReadable() string {
	return fmt.Sprintf("%s %s %s\n- %s=%s\n+ %s=%s",
		diff.Type,
		diff.MetaType,
		diff.Name,
		diff.Field,
		diff.OldValue,
		diff.Field,
		diff.NewValue)
}

func (diff FieldDiff) GetType() DiffType {
	return diff.Type
}

func (diff FieldDiff) GetName() string {
	return diff.Name
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

func GenerateModifiedSubDependencyDiff(pkgName string, parent string) DependencyDiff {
	return DependencyDiff{
		Type:     MODIFIED,
		MetaType: SUB_DEPENDENCY,
		Name:     pkgName,
		Parent:   parent,
	}
}

func GenerateAddedDependencyDiff(pkgName string, version string, parent string) DependencyDiff {
	return DependencyDiff{
		Type:     ADDED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Parent:   parent,
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
