package repository

import (
	"html/template"
	"os"
	"path"
	"reflect"
)

// genUpdate update.go
func (r *Repository) genUpdate(rt reflect.Type, _data *base) error {
	data := struct {
		Package     string
		GenQueryPkg string
		RepoPkg     string
		RepoPkgName string
		StructName  string
		Abbr        string
	}{
		Package:     _data.filename,
		GenQueryPkg: r.genQueryPkg,
		RepoPkg:     r.repoPkg,
		RepoPkgName: r.repoPkgName,
		StructName:  rt.Name(),
		Abbr:        _data.abbr,
	}

	file, err := os.Create(path.Join(_data.paths, "update.gen.go"))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	t, err := template.New(r.genUpdateTemplate()).Parse(r.genUpdateTemplate())
	if err != nil {
		return err
	}
	if err = t.Execute(file, data); err != nil {
		return err
	}
	return nil
}
