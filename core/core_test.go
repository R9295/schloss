package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDependencyFieldDiff(t *testing.T) {
	assert.Equal(t, GenerateDependencyFieldDiff("deno_core", "version", "old", "new"),
	FieldDiff{
		Type:     MODIFIED,
		MetaType: DEPENDENCY,
		Name:     "deno_core",
		Field:    "version",
		OldValue: "old",
		NewValue: "new",
	})
}

func TestGenerateModifiedSubDependencyDiff(t *testing.T) {
	assert.Equal(t, GenerateModifiedSubDependencyDiff("deno_core", "root_pkg"),
	DependencyDiff{
		Type:     MODIFIED,
		MetaType: SUB_DEPENDENCY,
		Name:     "deno_core",
		Parent:   "root_pkg",
	})
}

func TestGenerateAddedDependencyDiff(t *testing.T) {
	assert.Equal(t, GenerateAddedDependencyDiff("deno_core", "0.1", "root_pkg"),
	DependencyDiff{
		Type:     ADDED,
		MetaType: DEPENDENCY,
		Name:     "deno_core",
		Parent:   "root_pkg",
		Version:  "0.1",
	})
}

func TestGenerateRemovedDependencyDiff(t *testing.T) {
	assert.Equal(t, GenerateRemovedDependencyDiff("deno_core"),
	DependencyDiff{
		Type:     REMOVED,
		MetaType: DEPENDENCY,
		Name:     "deno_core",
		Parent:   "",
	})
}

func TestGenerateRemovedSubDependencyDiff(t *testing.T) {
	assert.Equal(t, GenerateRemovedSubDependencyDiff("assert", "deno_core"),
	DependencyDiff{
		Type:     REMOVED,
		MetaType: SUB_DEPENDENCY,
		Name:     "assert",
		Parent:   "deno_core",
	})
}

func TestGenerateAddedSubDependencyDiff(t *testing.T) {
	assert.Equal(t, GenerateAddedSubDependencyDiff("assert", "deno_core", "0.1"),
	DependencyDiff{
		Type:     ADDED,
		MetaType: SUB_DEPENDENCY,
		Name:     "assert",
		Parent:   "deno_core",
		Version:  "0.1",
	})
}
