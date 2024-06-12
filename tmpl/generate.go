package tmpl

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"path"

	"github.com/rogpeppe/go-internal/txtar"
)

//go:embed base.txt
var baseTemplate []byte

func Gen(base string) error {
	archive := txtar.Parse(baseTemplate)

	for _, file := range archive.Files {
		loc := path.Join(base, file.Name)
		if err := os.WriteFile(loc, file.Data, 0644); err != nil {
			return fmt.Errorf("write file %q failed: %w", file.Name, err)
		}
	}
	return nil
}

func GenHandler(structName string) ([]byte, error) {
	buf := new(bytes.Buffer)
	s := func(s string) {
		buf.WriteString(s)
		buf.WriteByte('\n')
	}
	s("package main")
	s("import (")
	s("\"net/http\"")
	s("\"github.com/go-json-experiment/json\"")
	s(")")
	s("func HandleCreate" + structName + "(w http.ResponseWriter, r *http.Request) {")
	s("var (")
	s("ctx = r.Context()")
	s(")")
	s("in := new(" + structName + ")")
	s("if err := json.UnmarshalRead(r.Body, &in); err != nil {")
	s("serveError(w, http.StatusBadRequest, err)")
	s("return")
	s("}")
	s("if err := db.Create" + structName + "(ctx, in); err != nil {")
	s("serveError(w, http.StatusInternalServerError, err)")
	s("return")
	s("}")
	s("serveJSON(w, http.StatusCreated, in)")
	s("}")
	s("")

	// handle list
	s("func HandleList" + structName + "(w http.ResponseWriter, r *http.Request) {")
	s("ctx := r.Context()")
	s("out, err := db.List" + structName + "(ctx, parseListQuery(r))")
	s("if err != nil {")
	s("serveError(w, http.StatusInternalServerError, err)")
	s("return")
	s("}")
	s("serveJSON(w, http.StatusOK, out)")
	s("}")

	return format.Source(buf.Bytes())
}

func GenDB(table string, structName string, columns []string, fields []string) ([]byte, error) {
	buf := new(bytes.Buffer)
	s := func(s string) {
		buf.WriteString(s)
		buf.WriteByte('\n')
	}

	s("package main")
	s("")
	s("import (")
	s("\"context\"")
	s(")")

	s("func (db *DB) Create" + structName + "(ctx context.Context, in *" + structName + ") error {")
	s("stmt, args, err := sq.Insert(\"" + table + "\").")

	s("Columns(")
	for _, c := range columns {
		fmt.Fprintf(buf, "\"%s\",", c)
	}
	s(").")

	s("Values(")
	for _, f := range fields {
		fmt.Fprintf(buf, "in.%s,", f)
	}
	s(").")
	s("Suffix(\"RETURNING *\").")
	s("ToSql()")
	s("if err != nil {")
	s("return err")
	s("}")
	s("return db.db.QueryRowxContext(ctx, stmt, args...).StructScan(in)")
	s("}")

	// list
	s("func (db *DB) List" + structName + "(ctx context.Context, params *ListParam) ([]" + structName + ", error) {")
	s("stmt, args, err := sq.Select(\"*\").")
	s("From(\"" + table + "\").")
	s("Limit(params.Limit).")
	s("Offset(params.Offset).")
	s("ToSql()")
	s("if err != nil {")
	s("return nil, err")
	s("}")

	s("out := make([]" + structName + ", 0)")
	s("err = db.db.SelectContext(ctx, &out, stmt, args...)")
	s("return out, err")
	s("}")
	return format.Source(buf.Bytes())
}
