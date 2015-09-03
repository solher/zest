package zest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCli runs tests on the zest Cli and SetCli methods.
func TestCli(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	zest := New()

	r.NotPanics(func() { _ = zest.Cli() })
	cli := zest.Cli()

	cli.Name = "foobar"
	zest.SetCli(cli)

	r.NotPanics(func() { _ = zest.Cli() })
	cli = zest.Cli()

	a.Equal(cli.Name, "foobar")
}
