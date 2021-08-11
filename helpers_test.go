package parser_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// loadStringFromFile loads a file and returns the entire contents as a string. Any
// leading or trailing whitespace is removed
func loadStringFromFile(t *testing.T, dir string) string {
	t.Helper()

	nameParts := strings.Split(t.Name(), "/")
	filename := fmt.Sprintf("%s/%s", dir, nameParts[len(nameParts)-1])
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		assert.Failf(t, "error in test setup", "unable to load file %s", filename)
	}

	return strings.TrimSpace(string(out))
}
