package test

import (
	"testing"
	"testing/quick"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestAbs(t *testing.T) {
	f := func(x int) bool {
		return Abs(x) >= 0
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
