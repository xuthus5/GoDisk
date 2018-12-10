/***********************

	页面渲染

************************/

package controllers

import (
	"GoDisk/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

// 网站首页  路由 /
func (this *MainController) Get() {
	this.TplName = "index.html"
}

// 后台首页  路由 /admin
func (this *MainController) Admin() {
	master := this.GetSession("master")
	if master == nil {
		this.Redirect("/login", 302)
	}
	this.Data["master"] = master
	this.Layout = "layout.html"
	this.TplName = "admin.html"
}

// 分类页面 路由 /admin/category
func (this *MainController) Category() {
	master := this.GetSession("master")
	if master == nil {
		this.Redirect("/login", 302)
	}
	this.Data["master"] = master
	this.Layout = "layout.html"
	this.TplName = "category.html"
}

// 分类修改页面 路由 /admin/category/update
func (this *MainController) CategoryUpdate() {
	master := this.GetSession("master")
	if master == nil {
		this.Redirect("/login", 302)
	}
	this.Data["master"] = master
	this.Data["category"] = models.GetOneCategoryInfo(this.GetString("id"))
	this.Layout = "layout.html"
	this.TplName = "category-update.html"
}

// 配置页面  路由 /admin/setting
func (this *MainController) Setting() {
	master := this.GetSession("master")
	if master == nil {
		this.Redirect("/login", 302)
	}
	this.Data["master"] = master
	this.Layout = "layout.html"
	this.TplName = "setting.html"
}

// 本地上传页面  路由 /admin/upload/local
func (this *MainController) LocalUpload() {
	master := this.GetSession("master")
	if master == nil {
		this.Redirect("/login", 302)
	}
	this.Data["master"] = master
	this.Layout = "layout.html"
	this.TplName = "attachment.html"
}

// 七牛云上传页面  路由 /admin/upload/qiniu
func (this *MainController) QiniuUpload() {
	master := this.GetSession("master")
	if master == nil {
		this.Redirect("/login", 302)
	}
	this.Data["master"] = master
	this.Layout = "layout.html"
	this.TplName = "qiniu-upload.html"
}

// 又拍云上传页面  路由 /admin/upload/upyun
func (this *MainController) UpyunUpload() {
	master := this.GetSession("master")
	if master == nil {
		this.Redirect("/login", 302)
	}
	this.Data["master"] = master
	this.Layout = "layout.html"
	this.TplName = "upyun-upload.html"
}

// 阿里云OSS上传页面  路由 /admin/upload/oss
func (this *MainController) OssUpload() {
	master := this.GetSession("master")
	if master == nil {
		this.Redirect("/login", 302)
	}
	this.Data["master"] = master
	this.Layout = "layout.html"
	this.TplName = "oss-upload.html"
}
