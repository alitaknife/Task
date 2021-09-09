package main

import (
	"fmt"
	"github.com/pkg/errors"
)

// 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，
// 是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

// 答：应该抛给上层，ErrNoRows 属于正常的情况，对应 rows == 0.理应抛给业务层对其进行处理

var ErrNoRows = errors.New("no rows")

// dao
func queryRow() error {
	return ErrNoRows
}

func main() {
	err := queryRow()
	if err == ErrNoRows {
		fmt.Printf("we find nothing: %v", err)
	}
}
