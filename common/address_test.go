package common

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestPackage(t *testing.T) { TestingT(t) }

type AddressSuite struct{}

var _ = Suite(&AddressSuite{})

func (s *AddressSuite) TestAddress(c *C) {
	_, err := NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6", ETHChain)
	c.Assert(err, NotNil)

	_, err = NewAddress("1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6", ETHChain)
	c.Check(err, NotNil)
	_, err = NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6X", ETHChain)
	c.Check(err, NotNil)
	_, err = NewAddress("bogus", ETHChain)
	c.Check(err, NotNil)
	c.Check(Address("").IsEmpty(), Equals, true)
	c.Check(NoAddress.Equals(Address("")), Equals, true)
	_, err = NewAddress("", ETHChain)
	c.Assert(err, NotNil)

	// eth tests
	_, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635a", ETHChain)
	c.Check(err, IsNil)
	// wrong length
	_, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635aaaaaaaaa", ETHChain)
	c.Check(err, NotNil)

	// good length but not valid hex string
	_, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab63zz", ETHChain)
	c.Check(err, NotNil)

}
