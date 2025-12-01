package repository

var keyWordMap = map[string]struct{}{
	"UnderlyingDB":    {},
	"UseDB":           {},
	"UseModel":        {},
	"UseTable":        {},
	"Quote":           {},
	"Debug":           {},
	"TableName":       {},
	"WithContext":     {},
	"As":              {},
	"Not":             {},
	"Or":              {},
	"Build":           {},
	"Columns":         {},
	"Hints":           {},
	"Distinct":        {},
	"Omit":            {},
	"Select":          {},
	"Where":           {},
	"Order":           {},
	"Group":           {},
	"Having":          {},
	"Limit":           {},
	"Offset":          {},
	"Join":            {},
	"LeftJoin":        {},
	"RightJoin":       {},
	"Save":            {},
	"Create":          {},
	"CreateInBatches": {},
	"Update":          {},
	"Updates":         {},
	"UpdateColumn":    {},
	"UpdateColumns":   {},
	"Find":            {},
	"FindInBatches":   {},
	"First":           {},
	"Take":            {},
	"Last":            {},
	"Pluck":           {},
	"Count":           {},
	"Scan":            {},
	"ScanRows":        {},
	"Row":             {},
	"Rows":            {},
	"Delete":          {},
	"Unscoped":        {},
	"Scopes":          {},
}

// escapeKeyword 转义关键字.
func escapeKeyword(keyword string) string {
	if _, ok := keyWordMap[keyword]; ok {
		return keyword + "_"
	}
	return keyword
}
