package main

import "testing"

func TestName(t *testing.T) {

	name := getName()

	if name != "World!" {
		t.Error("Respone from get Name is unexpected value")
	}
}
