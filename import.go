package repository

import (
	_ "github.com/gorm-gen/field"
	_ "github.com/gorm-gen/paginate"
	_ "github.com/gorm-gen/sharding"
	_ "github.com/opentracing/opentracing-go"
	_ "github.com/shopspring/decimal"
	_ "go.uber.org/zap"
	_ "gorm.io/gen"
	_ "gorm.io/gorm"
)
