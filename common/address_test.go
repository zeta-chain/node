package common

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

func TestAddress(t *testing.T) {
	addr := NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	require.EqualValuesf(t, NoAddress, addr, "address string should be empty")

	addr = NewAddress("bogus")
	require.EqualValuesf(t, NoAddress, addr, "address string should be empty")

	addr = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635a")
	require.EqualValuesf(t, "0x90f2b1ae50e6018230e90a33f98c7844a0ab635a", addr.String(), "address string should be equal")
}
