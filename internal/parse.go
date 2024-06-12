package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Parse extracts the first struct in the go file.
// it returns the map of fields and column name in the struct.
func Parse(filename string) (string, map[string]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return "", nil, err
	}

	// Iterate through the declarations in the file
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// Iterate through the specs in the declaration
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// Check if the type is a struct
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Extract the struct name
			structName := typeSpec.Name.Name

			// Extract the fields of the struct
			fieldMap := make(map[string]string)
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					dbTag := getTagValue(field.Tag.Value, "db")
					if dbTag != "" {
						fieldMap[name.Name] = dbTag
					}
				}
			}

			return structName, fieldMap, nil
		}
	}
	return "", nil, fmt.Errorf("no struct found in the file")
}

func getTagValue(tag, key string) string {
	tag = strings.Trim(tag, "`")
	tagParts := strings.Split(tag, " ")
	for _, part := range tagParts {
		if strings.HasPrefix(part, key+":") {
			return strings.TrimPrefix(part, key+":")
		}
	}
	return ""
}
