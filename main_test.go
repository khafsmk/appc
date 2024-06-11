package main

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestAll(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(e *testscript.Env) error {
			return nil
		},
	})
}
