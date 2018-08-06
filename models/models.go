package models

import (
		"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
		"log"
	"GoDisk/tools"
	)

type User struct {
	Id          int
	Username    string	`orm:"size(16)"`
	Password    string	`orm:"size(32)"`
	Created     string	`orm:"size(10)"`
}

type Classify struct {
	Id			int
	Label 		string	`orm:"size(64)"`
	Mark 		string  `orm:"size(64)"`
}

type File struct {
	Id			int
	Name 		string
	Mark 		string	`orm:"size(64)"`
	Path 		string
	Created 	string	`orm:"size(10)"`
}

//获取数据表数据
type Count struct {
	Num string
}

var dbc orm.Ormer
var dbx *sqlx.DB

func init() {
	// 注册驱动
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	// 注册默认数据库
	// 备注：此处第一个参数必须设置为“default”（因为我现在只有一个数据库），否则编译报错说：必须有一个注册DB的别名为 default
	orm.RegisterDataBase("default", "sqlite3", "static/db/data.db")
	// 需要在init中注册定义的model
	orm.RegisterModel(new(User),new(Classify),new(File))
	// 开启 orm 调试模式：开发过程中建议打开，release时需要关闭
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)

	dbc = orm.NewOrm()
	dbc.Using("default")
	//安装初始化
	CheckInstall()

	//sqlx
	dbx,_ = sqlx.Open("sqlite3","static/db/data.db")
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//defer dbx.Close()
}

// 检测是否初始化数据库
func CheckInstall(){
	user := &User{Username:"admin"}
	err := dbc.Read(user,"Username")
	if err != nil {
		user = &User{Username:"admin",Password:tools.StringToMd5("admin"),Created:tools.TimeToString()}
		classify := &Classify{Label:"默认",Mark:"default"}
		Register(user)
		AddClassify(classify)
	}
}

//注册
func Register(user *User) bool{
	_, err := dbc.Insert(user)
	if err == nil {
		return true
	}else {
		return false
	}
}

//登陆
func Login(user *User) (Code int64,Msg string){
	err := dbc.Read(user,"Username", "Password")
	if err == nil {
		return 1,"欢迎回来！"
	} else {
		return 0,"用户名或密码错误！"
	}
}

//添加分类
func AddClassify(info *Classify) bool {
	_,err := dbc.Insert(info)
	if err == nil{
		return true
	}else{
		return false
	}
}

//获取分类列表
func ApiClassifyList() *[]Classify{
		list := []Classify{}
	err := dbx.Select(&list, "select * from classify")
	if err != nil {
		log.Fatal(err.Error())
	}
	return &list
}

//获取文件列表
func ApiFileList() *[]File{
	list := []File{}
	err := dbx.Select(&list, "select * from File")
	if err != nil {
		log.Fatal(err.Error())
	}
	return &list
}

//文件上传 入数据库
func FileSave(info *File) bool {
	_,err := dbc.Insert(info)
	if err == nil{
		return true
	}else{
		return false
	}
}

//统计数据表的内容数量
func FindNumber(table string) *[]Count {
	num := []Count{}
	err := dbx.Select(&num,"select count(*) `num` from "+table)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &num
}