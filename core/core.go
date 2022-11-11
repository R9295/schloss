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
	FIELD          DiffMetaType = "field"
	META_FIELD     DiffMetaType = "meta-field"
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

type MetadataDiff struct {
	Type      DiffType     `json:"type"`
	MetaType  DiffMetaType `json:"meta_type"`
	FieldName string       `json:"field"`
	OldValue  string       `json:"old_value"`
	NewValue  string       `json:"new_value"`
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

func (diff MetadataDiff) GetName() string {
	return diff.FieldName
}

type AbsentFieldDiff struct {
	Type     DiffType     `json:"type"`
	MetaType DiffMetaType `json:"meta_type"`
	Name     string       `json:"dependency_name"`
	Field    string       `json:"field"`
}

func (diff AbsentFieldDiff) RenderHumanReadable() string {
	return fmt.Sprintf("%s %s %s of %s", diff.Type, diff.MetaType, diff.Field, diff.Name)
}

// TODO: not correct new value missing
func (diff MetadataDiff) RenderHumanReadable() string {
	return fmt.Sprintf("%s %s %s of %s", diff.Type, diff.MetaType, diff.FieldName, diff.OldValue)
}

func (diff AbsentFieldDiff) GetType() DiffType {
	return diff.Type
}

func (diff MetadataDiff) GetType() DiffType {
	return diff.Type
}

func (diff AbsentFieldDiff) GetName() string {
	return diff.Name
}

func MakeDependencyFieldDiff(
	pkgName string,
	fieldName string,
	oldVal string,
	newVal string,
) FieldDiff {
	return FieldDiff{
		Type:     MODIFIED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Field:    fieldName,
		OldValue: oldVal,
		NewValue: newVal,
	}
}

func MakeModifiedMetadataDiff(
	fieldName string,
	oldVal string,
	newVal string,
) MetadataDiff {
	return MetadataDiff{
		Type:      ADDED,
		MetaType:  META_FIELD,
		FieldName: fieldName,
		OldValue:  oldVal,
		NewValue:  newVal,
	}
}

func MakeModifiedSubDependencyDiff(pkgName string, parent string) DependencyDiff {
	return DependencyDiff{
		Type:     MODIFIED,
		MetaType: SUB_DEPENDENCY,
		Name:     pkgName,
		Parent:   parent,
	}
}

func MakeAddedDependencyDiff(pkgName string, version string, parent string) DependencyDiff {
	return DependencyDiff{
		Type:     ADDED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Parent:   parent,
		Version:  version,
	}
}

func MakeRemovedDependencyDiff(pkgName string) DependencyDiff {
	return DependencyDiff{
		Type:     REMOVED,
		MetaType: DEPENDENCY,
		Name:     pkgName,
		Parent:   "",
	}
}

func MakeRemovedSubDependencyDiff(pkgName, parent string) DependencyDiff {
	return DependencyDiff{
		Type:     REMOVED,
		MetaType: SUB_DEPENDENCY,
		Name:     pkgName,
		Parent:   parent,
	}
}

func MakeAddedSubDependencyDiff(pkgName string, parent string, version string) DependencyDiff {
	return DependencyDiff{
		Type:     ADDED,
		MetaType: SUB_DEPENDENCY,
		Name:     pkgName,
		Parent:   parent,
		Version:  version,
	}
}

func MakeAbsentFieldDiff(pkgName string, field string) AbsentFieldDiff {
	return AbsentFieldDiff{
		Type:     REMOVED,
		MetaType: FIELD,
		Name:     pkgName,
		Field:    field,
	}
}
