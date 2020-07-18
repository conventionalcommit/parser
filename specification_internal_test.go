package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeIsPatch(t *testing.T) {
	expectedTrue := []string{
		"FIX",
		"FIx",
		"FiX",
		"Fix",
		"fIX",
		"fIx",
		"fiX",
		"fix",
	}

	for _, ex := range expectedTrue {
		assert.True(t, typeIsPatch(ex, DefaultPatchTypes))
	}

	expectedFalse := []string{
		"",
		"feat",
	}

	for _, ex := range expectedFalse {
		assert.False(t, typeIsPatch(ex, DefaultPatchTypes))
	}
}
