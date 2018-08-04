package main

import (
	_ "GoDisk/routers"
	"github.com/astaxie/beego"
	_ "github.com/mattn/go-sqlite3"
		)

func main() {
	beego.SetStaticPath("/data", "data")
	beego.Run()
}
