package common

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestPackage(t *testing.T) { TestingT(t) }

type AddressSuite struct{}

var _ = Suite(&AddressSuite{})

func (s *AddressSuite) TestAddress(c *C) {
	addr, err := NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Assert(err, NotNil)

	_, err = NewAddress("1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(err, NotNil)
	_, err = NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6X")
	c.Check(err, NotNil)
	_, err = NewAddress("bogus")
	c.Check(err, NotNil)
	c.Check(Address("").IsEmpty(), Equals, true)
	c.Check(NoAddress.Equals(Address("")), Equals, true)
	_, err = NewAddress("")
	c.Assert(err, IsNil)

	// eth tests
	addr, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635a")
	c.Check(err, IsNil)
	c.Check(addr.IsChain(ETHChain), Equals, true)
	// wrong length
	_, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635aaaaaaaaa")
	c.Check(err, NotNil)

	// good length but not valid hex string
	_, err = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab63zz")
	c.Check(err, NotNil)

}
