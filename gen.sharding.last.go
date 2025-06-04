package repository

import (
	"html/template"
	"os"
	"path"
	"reflect"
)

// genCount count.go
func (r *Repository) genShardingLast(rt reflect.Type, _data *base) error {
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
	}

	file, err := os.Create(path.Join(_data.paths, "sharding.last.gen.go"))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	t, err := template.New(r.genShardingLastTemplate()).Parse(r.genShardingLastTemplate())
	if err != nil {
		return err
	}
	if err = t.Execute(file, data); err != nil {
		return err
	}
	return nil
}
