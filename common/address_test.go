package common

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type AddressSuite struct{}

var _ = Suite(&AddressSuite{})

func (s *AddressSuite) TestAddress(c *C) {
	_, err := NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6", GoerliChain())
	c.Assert(err, NotNil)

	_, err = NewAddress("1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6", GoerliChain())
	c.Check(err, NotNil)
	_, err = NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6X", GoerliChain())
	c.Check(err, NotNil)
	_, err = NewAddress("bogus", GoerliChain())
	c.Check(err, NotNil)
	c.Check(Address("").IsEmpty(), Equals, true)
	c.Check(NoAddress.Equals(Address("")), Equals, true)
	_, err = NewAddress("", GoerliChain())
	c.Assert(err, NotNil)

	// eth tests
	_, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635a", GoerliChain())
	c.Check(err, IsNil)
	// wrong length
	_, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635aaaaaaaaa", GoerliChain())
	c.Check(err, NotNil)

	// good length but not valid hex string
	_, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab63zz", GoerliChain())
	c.Check(err, NotNil)

}
