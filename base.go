package repository

type pkg struct {
	time    bool
	decimal bool // github.com/shopspring/decimal
	pf      bool // github.com/gorm-gen/plugin/field
	pfv     bool // github.com/gorm-gen/plugin/field/value
	reflect bool
	string  bool
}

type base struct {
	abbr                  string
	filename              string
	paths                 string
	modelName             string
	shardingKey           string
	shardingKeyType       string
	shardingKeyTypeFormat string
}
