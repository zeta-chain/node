package metaclient

import (
	"testing"
)

// Testing boilerplate... todo once implemented

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
}

var _ = Suite(&MySuite{})
