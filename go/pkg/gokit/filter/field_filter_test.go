package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldFilters_Compile(t *testing.T) {
	filters := &FieldFilters{}
	filters.AddFilter("a", "in(b, c)")
	filters.AddFilter("b", 12)
	filters.AddFilter("c", "gt(100)")

	filter := filters.Compile("d > 0 && d < 100")
	assert.NotEmpty(t, filter)
	assert.Equal(t, "a == in(b, c) && b == 12 && c == gt(100) && (d > 0 && d < 100)", string(filter))
}
