package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/khafsmk/appc/tmpl"
	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
)

var (
	genCmd = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Initialize a web project",
		Run: func(cmd *cobra.Command, args []string) {
			err := generate(args)
			cobra.CheckErr(err)
		},
	}

	fset = token.NewFileSet()
)

func generate(args []string) error {
	seedFile := args[0]
	if _, err := os.Stat(seedFile); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("directory %q does not exist", seedFile)
	}
	base := args[1]
	if _, err := os.Stat(base); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("init project %q does not exist", base)
	}

	// parse the seed.go file, extract the first struct,
	seed, err := parseFile(fset, seedFile)
	if err != nil {
		return fmt.Errorf("parse seed file: %w", err)
	}

	writeFile := func(name string, data []byte) error {
		f, err := os.Create(path.Join(base, name))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(data)
		return err
	}

	// generate the db
	b, err := tmpl.GenHandler(seed.StructName)
	if err != nil {
		return fmt.Errorf("generate handler: %w", err)
	}
	if err := writeFile(fmt.Sprintf("handle_%s.go", seed.TableName), b); err != nil {
		return fmt.Errorf("write handler gen: %w", err)
	}

	// generate the db
	b, err = tmpl.GenDB(seed.TableName, seed.StructName, seed.Fields, seed.Columns)
	if err != nil {
		return fmt.Errorf("generate db: %w", err)
	}
	if err := writeFile(fmt.Sprintf("db_%s.go", seed.TableName), b); err != nil {
		return fmt.Errorf("write db gen: %w", err)
	}

	// copy the seed file to the base
	if err := appendSeed(seed.Source, path.Join(base, "core.go")); err != nil {
		return fmt.Errorf("append seed file: %w", err)
	}

	return nil
}

func writeTo(w io.Writer, b []byte) {
	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}

func appendSeed(source []byte, modelFile string) error {
	f, err := os.OpenFile(modelFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("open core file: %w", err)
	}
	defer f.Close()

	if _, err = f.Write(source); err != nil {
		return err
	}

	// auto import missing lib
	b, err := imports.Process(modelFile, nil, nil)
	if err != nil {
		return err
	}
	err = os.WriteFile(modelFile, b, 0600)
	if err != nil {

		return err
	}

	return nil
}

const tableNamePrefix = "TableName: "

type SeedFile struct {
	Source     []byte
	TableName  string
	StructName string
	Fields     []string
	Columns    []string
}

func parseFile(fset *token.FileSet, filename string) (*SeedFile, error) {

	seed := &SeedFile{
		TableName:  "",
		StructName: "",
		Fields:     make([]string, 0),
		Columns:    make([]string, 0),
	}

	af, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// find struct name
	for _, decl := range af.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					seed.StructName = typeSpec.Name.Name
					break
				}
			}
		}
	}

	// find table name from comments
	for _, cmt := range af.Comments {
		text := cmt.Text()
		if strings.Contains(text, tableNamePrefix) {
			out := strings.Split(text, tableNamePrefix)
			table := strings.Trim(out[1], " ")
			seed.TableName = strings.Trim(table, "\n")
		}
	}

	// find columns and fields
	var structTypes []*ast.StructType
	ast.Inspect(af, func(n ast.Node) bool {
		if n, ok := n.(*ast.StructType); ok {
			structTypes = append(structTypes, n)
		}
		return true

	})

	for _, structType := range structTypes {

		for _, field := range structType.Fields.List {
			seed.Columns = append(seed.Columns, field.Names[0].Name)
			col := field.Tag.Value
			col = reflect.StructTag(col).Get("db")
			seed.Fields = append(seed.Fields, col)
		}
	}

	if b, err := extractSeedStruct(af, seed.StructName); err != nil {
		return nil, err
	} else {
		seed.Source = b
	}

	return seed, nil
}

func extractSeedStruct(af *ast.File, structName string) ([]byte, error) {
	// extract the struct
	// Inspect the AST to find the specified struct
	var structDecl *ast.GenDecl
	ast.Inspect(af, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if typeSpec.Name.Name == structName {
				structDecl = genDecl
				return false // stop the inspection
			}
		}
		return true
	})
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, structDecl)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
