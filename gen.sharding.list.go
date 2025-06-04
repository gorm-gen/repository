package repository

import (
	"html/template"
	"os"
	"path"
	"reflect"
)

// genCount count.go
func (r *Repository) genShardingList(rt reflect.Type, _data *base) error {
	data := struct {
		Package               string
		GenQueryPkg           string
		RepoPkg               string
		RepoPkgName           string
		StructName            string
		ModelPkg              string
		ModelName             string
		Abbr                  string
		ShardingKey           string
		ShardingKeyType       string
		ShardingKeyTypeFormat string
		ChanSign              template.HTML
		DecimalPkg            template.HTML
		ToShardingValue       template.HTML
		ShardingValueTo       string
	}{
		Package:               _data.filename,
		GenQueryPkg:           r.genQueryPkg,
		RepoPkg:               r.repoPkg,
		RepoPkgName:           r.repoPkgName,
		StructName:            rt.Name(),
		ModelPkg:              rt.PkgPath(),
		ModelName:             _data.modelName,
		Abbr:                  _data.abbr,
		ShardingKey:           _data.shardingKey,
		ShardingKeyType:       _data.shardingKeyType,
		ShardingKeyTypeFormat: _data.shardingKeyTypeFormat,
		ChanSign:              template.HTML("<-"),
		ToShardingValue:       "k",
		ShardingValueTo:       "shardingValue := v.ShardingValue",
	}

	if _data.shardingKeyType != "string" {
		_typeStart := _data.shardingKeyType + "("
		_typeEnd := ")"
		if _data.shardingKeyType == "int64" {
			_typeStart = ""
			_typeEnd = ""
		}
		data.DecimalPkg = `
	"github.com/shopspring/decimal"`
		data.ToShardingValue = `fmt.Sprintf("%d", k)`
		data.ShardingValueTo = `_shardingValue, _ := decimal.NewFromString(v.ShardingValue)
					shardingValue := ` + _typeStart + `_shardingValue.BigInt().Int64()` + _typeEnd
	}

	file, err := os.Create(path.Join(_data.paths, "sharding.list.gen.go"))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	t, err := template.New(r.genShardingListTemplate()).Parse(r.genShardingListTemplate())
	if err != nil {
		return err
	}
	if err = t.Execute(file, data); err != nil {
		return err
	}
	return nil
}
