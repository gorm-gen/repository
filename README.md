### [使用教程](https://github.com/gorm-gen/example)
```text
go get -u github.com/gorm-gen/repository
```
```go
package main

import (
	"fmt"

	"github.com/gorm-gen/repository"
	
	"demo/internal/models"
)

func main() {
	r := repository.New(
		repository.WithModule("demo"),
		repository.WithRepositoryPath("cmd/internal/repositories"),
		repository.WithGenQueryPkg("demo/internal/query"),
		repository.WithGormDBVar("global.DB"),
		repository.WithGormDBVarPkg("demo/internal/global"),
		repository.WithZapVar("global.Logger"),
		repository.WithZapVarPkg("demo/internal/global"),
	)
	err := r.Generate(
		models.User{},
		models.Admin{},
	)
	if err != nil {
		panic(err)
		return
	}
	err = r.ShardingGenerate(
		"Sharding",
		models.User{},
		models.Admin{},
	)
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("success")
}

```