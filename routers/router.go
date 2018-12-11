package routers

import (
	"GoDisk/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//页面路由
	beego.Router("/", &controllers.MainController{})                                          //网站首页
	beego.Router("/admin", &controllers.MainController{}, "*:Admin")                          //网站后台首页
	beego.Router("/admin/category", &controllers.MainController{}, "*:Category")              //上传文件分类管理
	beego.Router("/admin/category/update", &controllers.MainController{}, "*:CategoryUpdate") //上传文件分类管理
	beego.Router("/admin/setting", &controllers.MainController{}, "*:Setting")                //网站配置页面
	beego.Router("/admin/upload/local", &controllers.MainController{}, "*:LocalUpload")       //本地文件上传

	//用户模块
	beego.Router("/login", &controllers.UserController{}, "*:Login")     //用户登陆
	beego.Router("/logout", &controllers.UserController{}, "get:Logout") //用户注销登陆

	//接口Api
	beego.Router("/api/category/list", &controllers.ApiController{}, "get:CategoryList")      //分类列表
	beego.Router("/api/category/add", &controllers.ApiController{}, "post:CategoryAdd")       //分类添加
	beego.Router("/api/category/update", &controllers.ApiController{}, "post:CategoryUpdate") //分类修改
	beego.Router("/api/category/delete", &controllers.ApiController{}, "get:CategoryDelete")  //分类删除
	beego.Router("/api/site/config", &controllers.ApiController{}, "post:SiteConfig")         //网站配置
	beego.Router("/api/file/upload", &controllers.ApiController{}, "post:FileUpload")         //文件上传
	beego.Router("/api/file/list", &controllers.ApiController{}, "get:FileList")              //文件列表
	beego.Router("/api/file/delete", &controllers.ApiController{}, "get:FileDelete")          //文件删除

	//七牛云模块
	beego.Router("/admin/upload/qiniu", &controllers.MainController{}, "get:QiniuUpload")       //七牛云文件上传页面
	beego.Router("/api/upload/qiniu", &controllers.ApiController{}, "post:QiniuUpload")         //上传接口
	beego.Router("/api/file/qiniu/list", &controllers.ApiController{}, "get:QiniuList")         //七牛文件列表
	beego.Router("/api/file/qiniu/delete", &controllers.ApiController{}, "get:QiniuDeleteFile") //七牛文件删除

	//又拍云模块
	beego.Router("/admin/upload/upyun", &controllers.MainController{}, "get:UpyunUpload")       //又拍云上传页面
	beego.Router("/api/upload/upyun", &controllers.ApiController{}, "post:UpyunUpload")         //上传接口
	beego.Router("/api/file/upyun/list", &controllers.ApiController{}, "get:UpyunList")         //又拍云文件列表
	beego.Router("/api/file/upyun/delete", &controllers.ApiController{}, "get:UpyunDeleteFile") //又拍云文件删除

	//阿里云oss模块
	beego.Router("/admin/upload/oss", &controllers.MainController{}, "get:OssUpload")       //阿里云上传页面
	beego.Router("/api/upload/oss", &controllers.ApiController{}, "post:OssUpload")         //阿里云上传接口
	beego.Router("/api/file/oss/list", &controllers.ApiController{}, "get:OssList")         //阿里云文件列表
	beego.Router("/api/file/oss/delete", &controllers.ApiController{}, "get:OssDeleteFile") //阿里云文件删除

	//腾讯云cos模块
	beego.Router("/admin/upload/cos", &controllers.MainController{}, "get:CosUpload")       //腾讯云上传页面
	beego.Router("/api/upload/cos", &controllers.ApiController{}, "post:CosUpload")         //腾讯云上传接口
	beego.Router("/api/file/cos/list", &controllers.ApiController{}, "get:CosList")         //腾讯云文件列表
	beego.Router("/api/file/cos/delete", &controllers.ApiController{}, "get:CosDeleteFile") //腾讯云文件删除
}
