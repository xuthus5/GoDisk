package controllers

import (
	"GoDisk/models"
	"GoDisk/tools"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.TplName = "index.html"
}

func (this *MainController) Admin(){
	sess := this.GetSession("Username")
	if sess == nil{
		this.Redirect("/login",302)
	}
	this.Data["Username"] = sess
	this.Data["File"] = models.FindNumber("file")
	this.Data["Classify"] = models.FindNumber("classify")
	this.Layout = "layout.html"
	this.TplName = "admin.html"
}

func (this *MainController) Classify() {
	sess := this.GetSession("Username")
	if sess == nil{
		this.Redirect("/login",302)
	}
	label := this.GetString("label")
	mark := this.GetString("mark")
	if label == "" || mark == "" {
		this.Data["Username"] = sess
		this.Data["list"] = models.ApiClassifyList()
		this.Layout = "layout.html"
		this.TplName = "classify.html"
	}else{
		info := &models.Classify{Label:label,Mark:mark}
		code := models.AddClassify(info)
		var data *ResultData
		if code == false {
			data = &ResultData{Code:0,Title:"结果:",Msg:"操作失败！"}
		}else{
			tools.DirCreate("data/"+mark)
			data = &ResultData{Code:1,Title:"结果:",Msg:"操作成功！"}
		}
		this.Data["json"] = data
		this.ServeJSON()
	}
}

func (this *MainController) Setting() {
	sess := this.GetSession("Username")
	if sess == nil{
		this.Redirect("/login",302)
	}
	this.Data["Username"] = sess
	this.Data["Qiniu"] = models.SiteConfigMap()
	this.Layout = "layout.html"
	this.TplName = "setting.html"
}

func (this *MainController) PostSetting() {
	config := &models.QiniuConfig{}
	if err := this.ParseForm(config); err != nil {
		data := &ResultData{Code:0,Title:"结果:",Msg:"数据更新失败！"}
		this.Data["json"] = data
		this.ServeJSON()
	}else{
		models.SiteConfig(*config)
		data := &ResultData{Code:1,Title:"结果:",Msg:"数据更新成功！"}
		this.Data["json"] = data
		this.ServeJSON()
	}
}

func (this *MainController) LocalUpload() {
	sess := this.GetSession("Username")
	if sess == nil{
		this.Redirect("/login",302)
	}
	this.Data["Username"] = sess
	this.Data["classify"] = models.ApiClassifyList()
	this.Data["list"] = models.ApiFileList()
	this.Layout = "layout.html"
	this.TplName = "localUpload.html"
}