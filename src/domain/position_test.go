package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPosition(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var x int = 1
		var y int = 2

		pos := NewPosition(x, y)

		assert.Equal(t, x, pos.x)
		assert.Equal(t, y, pos.y)
	})
}

func TestPositionX(t *testing.T) {
	t.Parallel()

	var x int = 1

	pos := NewPosition(x, 0)

	assert.Equal(t, x, pos.X())
}

func TestPositionY(t *testing.T) {
	t.Parallel()

	var y int = 1

	pos := NewPosition(0, y)

	assert.Equal(t, y, pos.Y())
}
