package mapwatcher

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestSubTree(t *testing.T) {
	treeRoot := NewSubTreeRoot()

	treeRoot.Add("long.addresses", "addr1")
	treeRoot.Add("long.addresses", "addr2")
	treeRoot.Add("long.addresses", "addr3")
	treeRoot.Remove("long.addresses", "addr1")

	vals := treeRoot.Get("long")
	spew.Dump(vals)
}

func TestSliceSwap(t *testing.T) {
	slices := []string{"1", "2", "3", "4"}

	slices[0], slices[2] = slices[2], "0"

	spew.Dump(slices)

	assert.Equal(t, "3", slices[0])
	assert.Equal(t, "0", slices[2])
}
