package repository

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"reflect"
)

type Condition string

type Update string

type Order string

func (r *Repository) genBase(rt reflect.Type, _data *base) error {
	var timePkg, decimalPkg, pfPkg, pfvPkg, reflectPkg, stringPkg bool

	conditions, _pkg := r.genConditionOpt(rt, _data)
	if _pkg.time {
		timePkg = true
	}
	if _pkg.decimal {
		decimalPkg = true
	}
	if _pkg.pf {
		pfPkg = true
	}
	if _pkg.pfv {
		pfvPkg = true
	}
	if _pkg.reflect {
		reflectPkg = true
	}
	if _pkg.string {
		stringPkg = true
	}

	_conditions := make([]template.HTML, 0, len(conditions))
	for _, condition := range conditions {
		_conditions = append(_conditions, template.HTML(condition))
	}

	updates, _pkg := r.genUpdateOpt(rt, _data)
	if _pkg.time {
		timePkg = true
	}
	if _pkg.decimal {
		decimalPkg = true
	}
	if _pkg.pf {
		pfPkg = true
	}
	if _pkg.pfv {
		pfvPkg = true
	}

	_updates := make([]template.HTML, 0, len(updates))
	for _, update := range updates {
		_updates = append(_updates, template.HTML(update))
	}

	var imports []template.HTML
	var wrap bool
	if reflectPkg {
		imports = append(imports, `	"reflect"
`)
		wrap = true
	}
	if stringPkg {
		imports = append(imports, `	"strings"
`)
		wrap = true
	}
	if timePkg {
		imports = append(imports, `	"time"
`)
		wrap = true
	}
	if wrap {
		imports = append(imports, "\n")
	}

	if pfPkg {
		imports = append(imports, `	f "github.com/gorm-gen/field"
`)
	}
	if pfvPkg {
		imports = append(imports, `	"github.com/gorm-gen/field/value"
`)
	}
	if decimalPkg {
		imports = append(imports, `	"github.com/shopspring/decimal"
`)
	}

	if r.zapVarPkg == r.gormDBVarPkg {
		imports = append(imports, template.HTML(fmt.Sprintf(`	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"%s"

	"%s"

	"%s"`, r.zapVarPkg, r.genQueryPkg, r.repoPkg)))
	}

	if r.zapVarPkg != r.gormDBVarPkg {
		imports = append(imports, template.HTML(fmt.Sprintf(`	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"%s"

	"%s"

	"%s"

	"%s"`, r.zapVarPkg, r.gormDBVarPkg, r.genQueryPkg, r.repoPkg)))
	}

	data := struct {
		Package     string
		ZapVarPkg   string
		GenQueryPkg string
		RepoPkg     string
		RepoPkgName string
		StructName  string
		Abbr        string
		GormDBVar   string
		ZapVar      string
		Imports     []template.HTML
		Conditions  []template.HTML
		Updates     []template.HTML
		Orders      []Order
	}{
		Package:     _data.filename,
		ZapVarPkg:   r.zapVarPkg,
		GenQueryPkg: r.genQueryPkg,
		RepoPkg:     r.repoPkg,
		RepoPkgName: r.repoPkgName,
		StructName:  rt.Name(),
		Abbr:        _data.abbr,
		GormDBVar:   r.gormDBVar,
		ZapVar:      r.zapVar,
		Imports:     imports,
		Conditions:  _conditions,
		Updates:     _updates,
		Orders:      r.genOrderOpt(rt, _data),
	}

	file, err := os.Create(path.Join(_data.paths, "base.gen.go"))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	t, err := template.New(r.genBaseTemplate()).Parse(r.genBaseTemplate())
	if err != nil {
		return err
	}
	if err = t.Execute(file, data); err != nil {
		return err
	}
	return nil
}
