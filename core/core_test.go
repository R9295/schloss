package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeDependencyFieldDiff(t *testing.T) {
	assert.Equal(t, MakeDependencyFieldDiff("deno_core", "version", "old", "new"),
		FieldDiff{
			Type:     MODIFIED,
			MetaType: DEPENDENCY,
			Name:     "deno_core",
			Field:    "version",
			OldValue: "old",
			NewValue: "new",
		})
}

func TestMakeModifiedSubDependencyDiff(t *testing.T) {
	assert.Equal(t, MakeModifiedSubDependencyDiff("deno_core", "root_pkg"),
		DependencyDiff{
			Type:     MODIFIED,
			MetaType: SUB_DEPENDENCY,
			Name:     "deno_core",
			Parent:   "root_pkg",
		})
}

func TestMakeAddedDependencyDiff(t *testing.T) {
	assert.Equal(t, MakeAddedDependencyDiff("deno_core", "0.1", "root_pkg"),
		DependencyDiff{
			Type:     ADDED,
			MetaType: DEPENDENCY,
			Name:     "deno_core",
			Parent:   "root_pkg",
			Version:  "0.1",
		})
}

func TestMakeRemovedDependencyDiff(t *testing.T) {
	assert.Equal(t, MakeRemovedDependencyDiff("deno_core"),
		DependencyDiff{
			Type:     REMOVED,
			MetaType: DEPENDENCY,
			Name:     "deno_core",
			Parent:   "",
		})
}

func TestMakeRemovedSubDependencyDiff(t *testing.T) {
	assert.Equal(t, MakeRemovedSubDependencyDiff("assert", "deno_core"),
		DependencyDiff{
			Type:     REMOVED,
			MetaType: SUB_DEPENDENCY,
			Name:     "assert",
			Parent:   "deno_core",
		})
}

func TestMakeAddedSubDependencyDiff(t *testing.T) {
	assert.Equal(t, MakeAddedSubDependencyDiff("assert", "deno_core", "0.1"),
		DependencyDiff{
			Type:     ADDED,
			MetaType: SUB_DEPENDENCY,
			Name:     "assert",
			Parent:   "deno_core",
			Version:  "0.1",
		})
}
