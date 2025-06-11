package repository

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
)

type Repository struct {
	module       string
	repoPath     string
	repoPkg      string
	repoPkgName  string
	genQueryPkg  string
	gormDBVar    string
	gormDBVarPkg string
	zapVar       string
	zapVarPkg    string
}

func New(opts ...Option) *Repository {
	repo := &Repository{
		module:       "demo",
		repoPath:     "internal/repositories",
		repoPkgName:  "repositories",
		genQueryPkg:  "demo/internal/query",
		gormDBVar:    "global.DB",
		gormDBVarPkg: "demo/internal/global",
		zapVar:       "global.Logger",
		zapVarPkg:    "demo/internal/global",
	}

	for _, opt := range opts {
		opt(repo)
	}

	repoPathArr := strings.Split(repo.repoPath, "/")
	repo.repoPkgName = repoPathArr[len(repoPathArr)-1]
	repo.repoPkg = path.Join(repo.module, repo.repoPath)

	return repo
}

// Generate 生成普通仓库
func (r *Repository) Generate(models ...interface{}) error {
	if len(models) == 0 {
		return nil
	}

	if err := r.generateRepositoriesBase(); err != nil {
		return err
	}

	for _, model := range models {
		if err := r.generate(model, ""); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) generateRepositoriesBase() error {
	if err := os.MkdirAll(r.repoPath, os.ModePerm); err != nil {
		return err
	}

	// repositories/base.go
	if err := r.repositoriesBase(); err != nil {
		return err
	}

	return nil
}

// ShardingGenerate 生成分表仓库
func (r *Repository) ShardingGenerate(shardingStructName string, models ...interface{}) error {
	if len(models) == 0 {
		return nil
	}

	if err := r.generateRepositoriesBase(); err != nil {
		return err
	}

	for _, model := range models {
		if err := r.generate(model, shardingStructName); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) generate(model interface{}, shardingStructName string) error {
	rt := reflect.TypeOf(model)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	abbr := strings.ToLower(rt.Name()[:1])
	filename := abbr + rt.Name()[1:]
	paths := path.Join(r.repoPath, filename)
	if err := os.MkdirAll(paths, os.ModePerm); err != nil {
		return err
	}
	modelPkgArr := strings.Split(rt.PkgPath(), "/")
	modelName := modelPkgArr[len(modelPkgArr)-1]

	_base := &base{
		abbr:      abbr,
		filename:  filename,
		paths:     paths,
		modelName: modelName,
	}

	// base.go
	if err := r.genBase(rt, _base); err != nil {
		return err
	}

	// count.go
	if err := r.genCount(rt, _base); err != nil {
		return err
	}

	// create.go
	if err := r.genCreate(rt, _base); err != nil {
		return err
	}

	// delete.go
	if err := r.genDelete(rt, _base); err != nil {
		return err
	}

	// first.go
	if err := r.genFirst(rt, _base); err != nil {
		return err
	}

	// last.go
	if err := r.genLast(rt, _base); err != nil {
		return err
	}

	// list.go
	if err := r.genList(rt, _base); err != nil {
		return err
	}

	// pluck.go
	if err := r.genPluck(rt, _base); err != nil {
		return err
	}

	// scan.go
	if err := r.genScan(rt, _base); err != nil {
		return err
	}

	// sum.go
	if err := r.genSum(rt, _base); err != nil {
		return err
	}

	// take.go
	if err := r.genTake(rt, _base); err != nil {
		return err
	}

	// update.go
	if err := r.genUpdate(rt, _base); err != nil {
		return err
	}

	if shardingStructName == "" {
		return nil
	}

	shardingKeyExist := false
	var shardingKeyType string
	var shardingKeyTypeFormat string

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Name != shardingStructName {
			continue
		}
		allowType := false
		typ := field.Type.String()
		if r.isInt(typ) {
			allowType = true
			shardingKeyTypeFormat = "d"
		}
		if r.isString(typ) {
			allowType = true
			shardingKeyTypeFormat = "s"
		}
		if !allowType {
			return fmt.Errorf("%s sharding key %s type %s not support", rt.Name(), shardingStructName, typ)
		}
		shardingKeyType = strings.Trim(typ, "*")
		shardingKeyExist = true
		break
	}

	if !shardingKeyExist {
		return fmt.Errorf("%s not exist sharding key %s", rt.Name(), shardingStructName)
	}

	_base.shardingKey = shardingStructName
	_base.shardingKeyType = shardingKeyType
	_base.shardingKeyTypeFormat = shardingKeyTypeFormat

	// sharding.count.go
	if err := r.genShardingCount(rt, _base); err != nil {
		return err
	}

	// sharding.create.go
	if err := r.genShardingCreate(rt, _base); err != nil {
		return err
	}

	// sharding.delete.go
	if err := r.genShardingDelete(rt, _base); err != nil {
		return err
	}

	// sharding.first.go
	if err := r.genShardingFirst(rt, _base); err != nil {
		return err
	}

	// sharding.last.go
	if err := r.genShardingLast(rt, _base); err != nil {
		return err
	}

	// sharding.list.go
	if err := r.genShardingList(rt, _base); err != nil {
		return err
	}

	// sharding.sum.go
	if err := r.genShardingSum(rt, _base); err != nil {
		return err
	}

	// sharding.take.go
	if err := r.genShardingTake(rt, _base); err != nil {
		return err
	}

	// sharding.update.go
	if err := r.genShardingUpdate(rt, _base); err != nil {
		return err
	}

	return nil
}
