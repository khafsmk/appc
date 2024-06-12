package tmpl

import (
	"os"
	"testing"
)

func TestGenDB(t *testing.T) {

	f, err := os.CreateTemp("", "gendb")
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()
	buf, err := GenDB("users", "User", []string{"id", "name"}, []string{"ID", "Name"})
	if err != nil {
		t.Fatal(err)
	}
	f.Write(buf)

	t.Log(f.Name())
}

func TestGenHandler(t *testing.T) {
	f, err := os.CreateTemp("", "gendb")
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()
	buf, err := GenHandler("User")
	if err != nil {
		t.Fatal(err)
	}
	f.Write(buf)

	t.Log(f.Name())
}
