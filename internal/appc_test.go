package appc_test

import (
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestAll(t *testing.T) {

	files, err := filepath.Glob("testdata/*.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(strings.TrimSuffix(filepath.Base(file), ".txt"), func(t *testing.T) {
			a, err := txtar.ParseFile(file)
			if err != nil {
				t.Fatal(err)
			}
			for i := 0; i < len(a.Files); i++ {

			}
		})
	}

}
