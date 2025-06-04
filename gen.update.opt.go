package repository

import (
	"fmt"
	"reflect"
	"strings"
)

func (r *Repository) intUpdate(fieldName, fieldType string, rt reflect.Type, _data *base) []Update {
	var updates []Update

	update := fmt.Sprintf(`
// Update%[1]sAdd +=
func Update%[1]sAdd(v ...%[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		_v := %[2]s(1)
		if len(v) > 0 {
			_v = v[0]
		}
		if %[3]s.newTableName != nil {
			return %[3]s.q.%[4]s.Table(*%[3]s.newTableName).%[1]s.Add(_v)
		}
		return %[3]s.q.%[4]s.%[1]s.Add(_v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	update = fmt.Sprintf(`
// Update%[1]sSub -=
func Update%[1]sSub(v ...%[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		_v := %[2]s(1)
		if len(v) > 0 {
			_v = v[0]
		}
		if %[3]s.newTableName != nil {
			return %[3]s.q.%[4]s.Table(*%[3]s.newTableName).%[1]s.Sub(_v)
		}
		return %[3]s.q.%[4]s.%[1]s.Sub(_v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	update = fmt.Sprintf(`
// Update%[1]sMul *=
func Update%[1]sMul(v %[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		if %[3]s.newTableName != nil {
			return %[3]s.q.%[4]s.Table(*%[3]s.newTableName).%[1]s.Mul(v)
		}
		return %[3]s.q.%[4]s.%[1]s.Mul(v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	update = fmt.Sprintf(`
// Update%[1]sDiv /=
func Update%[1]sDiv(v %[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		if %[3]s.newTableName != nil {
			return %[3]s.q.%[4]s.Table(*%[3]s.newTableName).%[1]s.Div(v)
		}
		return %[3]s.q.%[4]s.%[1]s.Div(v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	return updates
}

func (r *Repository) decimalUpdate(fieldName, fieldType string, rt reflect.Type, _data *base) []Update {
	var updates []Update

	update := fmt.Sprintf(`
func Update%[1]s(v %[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		if %[3]s.newTableName != nil {
			return f.NewDecimal(%[3]s.q.%[4]s.%[1]s, f.WithTableName(*%[3]s.newTableName)).Value(v)
		}
		return f.NewDecimal(%[3]s.q.%[4]s.%[1]s).Value(v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	update = fmt.Sprintf(`
// Update%[1]sAdd +=
func Update%[1]sAdd(v %[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		if %[3]s.newTableName != nil {
			return f.NewDecimal(%[3]s.q.%[4]s.%[1]s, f.WithTableName(*%[3]s.newTableName)).Add(v)
		}
		return f.NewDecimal(%[3]s.q.%[4]s.%[1]s).Add(v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	update = fmt.Sprintf(`
// Update%[1]sSub -=
func Update%[1]sSub(v %[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		if %[3]s.newTableName != nil {
			return f.NewDecimal(%[3]s.q.%[4]s.%[1]s, f.WithTableName(*%[3]s.newTableName)).Sub(v)
		}
		return f.NewDecimal(%[3]s.q.%[4]s.%[1]s).Sub(v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	update = fmt.Sprintf(`
// Update%[1]sMul *=
func Update%[1]sMul(v %[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		if %[3]s.newTableName != nil {
			return f.NewDecimal(%[3]s.q.%[4]s.%[1]s, f.WithTableName(*%[3]s.newTableName)).Mul(v)
		}
		return f.NewDecimal(%[3]s.q.%[4]s.%[1]s).Mul(v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	update = fmt.Sprintf(`
// Update%[1]sDiv /=
func Update%[1]sDiv(v %[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		if %[3]s.newTableName != nil {
			return f.NewDecimal(%[3]s.q.%[4]s.%[1]s, f.WithTableName(*%[3]s.newTableName)).Div(v)
		}
		return f.NewDecimal(%[3]s.q.%[4]s.%[1]s).Div(v)
	}
}
`, fieldName, fieldType, _data.abbr, rt.Name())
	updates = append(updates, Update(update))

	return updates
}

func (r *Repository) genUpdateOpt(rt reflect.Type, _data *base) (updates []Update, _pkg *pkg) {
	_pkg = &pkg{}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if strings.ToLower(field.Name) == "id" {
			continue
		}
		typ := field.Type.String()
		if !r.allowType(typ) {
			continue
		}
		fieldType := strings.Trim(field.Type.String(), "*")

		if !r.isDecimal(typ) {
			update := fmt.Sprintf(`
func Update%[1]s(v %[2]s) UpdateOption {
	return func(%[3]s *%[4]s) field.AssignExpr {
		if %[3]s.newTableName != nil {
			return %[3]s.q.%[4]s.Table(*%[3]s.newTableName).%[1]s.Value(v)
		}
		return %[3]s.q.%[4]s.%[1]s.Value(v)
	}
}
`, field.Name, fieldType, _data.abbr, rt.Name())
			updates = append(updates, Update(update))
		}

		if r.isTime(typ) {
			_pkg.time = true
		}

		if r.isInt(typ) {
			updates = append(updates, r.intUpdate(field.Name, fieldType, rt, _data)...)
		}
		if r.isDecimal(typ) {
			_pkg.decimal = true
			_pkg.pf = true
			_pkg.pfv = true
			updates = append(updates, r.decimalUpdate(field.Name, fieldType, rt, _data)...)
		}

		if !strings.Contains(typ, "*") {
			continue
		}

		update := fmt.Sprintf(`
// Update%[1]sNull set null
func Update%[1]sNull() UpdateOption {
	return func(%[2]s *%[3]s) field.AssignExpr {
		if %[2]s.newTableName != nil {
			return %[2]s.q.%[3]s.Table(*%[2]s.newTableName).%[1]s.Null()
		}
		return %[2]s.q.%[3]s.%[1]s.Null()
	}
}
`, field.Name, _data.abbr, rt.Name())
		updates = append(updates, Update(update))
	}

	return
}
