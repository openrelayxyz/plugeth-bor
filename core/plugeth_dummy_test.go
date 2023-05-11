package core

import (
	"testing"
)

func TestAlwaysPasses(t *testing.T) {
	if true != true {
		t.Errorf("Expected true, but got false")
	}
}