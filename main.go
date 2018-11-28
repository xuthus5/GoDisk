package main

import (
	_ "GoDisk/routers"
	"github.com/astaxie/beego"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 自定义方法
	_ = beego.AddFuncMap("sc", SiteConfig) //调取网站配置 直接抽调数据库配置字段

	//自定义开放性资源路径
	beego.SetStaticPath("/file", "file")
	beego.Run()
}