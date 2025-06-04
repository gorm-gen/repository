package repository

import (
	"html/template"
	"os"
	"path"
	"reflect"
)

// genCount count.go
func (r *Repository) genShardingDelete(rt reflect.Type, _data *base) error {
	data := struct {
		Package               string
		GenQueryPkg           string
		RepoPkg               string
		RepoPkgName           string
		StructName            string
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
		Abbr:                  _data.abbr,
		ShardingKey:           _data.shardingKey,
		ShardingKeyType:       _data.shardingKeyType,
		ShardingKeyTypeFormat: _data.shardingKeyTypeFormat,
		ChanSign:              template.HTML("<-"),
	}

	file, err := os.Create(path.Join(_data.paths, "sharding.delete.gen.go"))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	t, err := template.New(r.genShardingDeleteTemplate()).Parse(r.genShardingDeleteTemplate())
	if err != nil {
		return err
	}
	if err = t.Execute(file, data); err != nil {
		return err
	}
	return nil
}
