/************************

	SQLite配置

*************************/

package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
)

var dbc orm.Ormer
var dbx *sqlx.DB

func init() {
	// 注册驱动
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	// 注册默认数据库
	orm.RegisterDataBase("default", "sqlite3", "data.db")
	// 需要在init中注册定义的model
	orm.RegisterModel(new(Category), new(Attachment), new(Config))
	// 开启 orm 调试模式：开发过程中建议打开，release时需要关闭
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)

	dbc = orm.NewOrm()
	dbc.Using("default")

	//检测是否安装
	isInstall := dbc.Read(&Config{Option: "IsInstall", Value: "yes"}, "Option", "Value")
	if isInstall != nil {
		Initialization()
	}
	//sqlx
	dbx, _ = sqlx.Open("sqlite3", "data.db")
}

//安装初始化
func Initialization() {
	//配置表初始化
	dbc.Insert(&Config{Option: "IsInstall", Value: "yes", Addition: ""})
	dbc.Insert(&Config{Option: "WebTitle", Value: "GoDisk", Addition: ""})
	dbc.Insert(&Config{Option: "Author", Value: "admin", Addition: ""})
	dbc.Insert(&Config{Option: "Password", Value: "admin", Addition: ""})
	dbc.Insert(&Config{Option: "CopyRight", Value: "GoDisk", Addition: ""})
	dbc.Insert(&Config{Option: "LogoUrl", Value: "/static/images/user-head-image.jpeg", Addition: ""})
	dbc.Insert(&Config{Option: "Keywords", Value: "", Addition: ""})
	dbc.Insert(&Config{Option: "Description", Value: "", Addition: ""})
	//默认七牛云配置信息写入数据库
	qiniuConfig := QiniuConfigOption{Zone: "storage.ZoneHuabei"}
	t := reflect.TypeOf(qiniuConfig)
	v := reflect.ValueOf(qiniuConfig)
	for i := 0; i < t.NumField(); i++ {
		dbc.Insert(&Config{Option: t.Field(i).Name, Value: v.Field(i).String(), Addition: "Qiniu"})
	}

	//默认又拍云配置信息写入数据库
	upyunConfig := UpyunConfigOption{}
	t = reflect.TypeOf(upyunConfig)
	v = reflect.ValueOf(upyunConfig)
	for i := 0; i < t.NumField(); i++ {
		dbc.Insert(&Config{Option: t.Field(i).Name, Value: v.Field(i).String(), Addition: "Upyun"})
	}

	//默认阿里云配置信息写入数据库
	OssConfig := OssConfigOption{}
	t = reflect.TypeOf(OssConfig)
	v = reflect.ValueOf(OssConfig)
	for i := 0; i < t.NumField(); i++ {
		dbc.Insert(&Config{Option: t.Field(i).Name, Value: v.Field(i).String(), Addition: "OSS"})
	}

	//默认腾讯云配置信息写入数据库
	CosConfig := CosConfigOption{}
	t = reflect.TypeOf(CosConfig)
	v = reflect.ValueOf(CosConfig)
	for i := 0; i < t.NumField(); i++ {
		dbc.Insert(&Config{Option: t.Field(i).Name, Value: v.Field(i).String(), Addition: "COS"})
	}

	//分类表初始化
	dbc.Insert(&Category{
		Name:        "默认",
		Key:         "default",
		Description: "默认的分类",
	})
}
