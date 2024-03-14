package attrvalue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEquals(t *testing.T) {
	t.Run("same elements", func(t *testing.T) {
		set1 := newSet([]int{1, 2, 3})
		set2 := newSet([]int{1, 2, 3})
		assert.True(t, set1.equals(set2))
	})

	t.Run("different elements", func(t *testing.T) {
		set1 := newSet([]int{1, 2, 3})
		set2 := newSet([]int{4, 5, 6})
		assert.False(t, set1.equals(set2))
	})

	t.Run("same elements different order", func(t *testing.T) {
		set1 := newSet([]int{1, 2, 3})
		set2 := newSet([]int{3, 2, 1})
		assert.True(t, set1.equals(set2))
	})

	t.Run("elements with different length", func(t *testing.T) {
		set1 := newSet([]int{1, 2, 3, 4})
		set2 := newSet([]int{3, 2, 1})
		assert.False(t, set1.equals(set2))
	})

	t.Run("self comparison", func(t *testing.T) {
		set1 := newSet([]int{1, 2, 3})
		assert.True(t, set1.equals(set1))
	})
}
