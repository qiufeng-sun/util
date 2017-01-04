package scribe

import (
	"testing"
)

//
func TestScribeNotInit(t *testing.T) {
	Log("test", "test not init scribe")
}

//
func TestScribe(t *testing.T) {
	InitScribe("localhost:7915", -1, -1)
	defer CloseScribe()

	Log("test", "send to scribe")
}
