package main

import (
	"fmt"

	"github.com/rogpeppe/go-internal/txtar"
)

func main() {
	// cmd.Execute()
	archive, err := txtar.ParseFile("./tpl/tpl.txt")
	if err != nil {
		panic(err)
	}

	for _, file := range archive.Files {
		fmt.Printf("file.Name: %v\n", file.Name)
		fmt.Printf("file.Data: %v\n", string(file.Data))
	}
}
