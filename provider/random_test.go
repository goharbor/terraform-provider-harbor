package provider

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	result := randomString(15)
	length := len(result)

	if length != 15 {
		t.Error("Failed not the correct length")
	}
}
