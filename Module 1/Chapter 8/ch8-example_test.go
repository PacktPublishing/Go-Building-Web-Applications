package example_test

import (
	"testing"
)

func TestSquare(t *Testing.T) {
	if Square(4) != 16 {
		t.Error("expected", 16, "got", 4)
	}
}
