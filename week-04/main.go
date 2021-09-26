package main

import (
	_ "week-04/boot"
	_ "week-04/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
