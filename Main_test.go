package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// t.Run("Check failed connection", func(t *testing.T) {
	// 	port := 12234

	// 	conn, err := OpenConnectionGRPC(port)
	// 	assert.Error(t, err)
	// 	assert.Empty(t, conn)
	// })

	t.Run("Check failed connection empty params", func(t *testing.T) {
		port := 0

		conn, err := OpenConnectionGRPC(port)
		assert.Error(t, err)
		assert.Empty(t, conn)
	})
	t.Run("Check Success connection", func(t *testing.T) {
		port := 6000

		conn, err := OpenConnectionGRPC(port)
		assert.NoError(t, err)
		assert.NotEmpty(t, conn)
	})
}
