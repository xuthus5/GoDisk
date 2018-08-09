package controllers

import (
	"github.com/astaxie/beego"
	"GoDisk/models"
	"GoDisk/tools"
	"regexp"
	"strings"
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
	if err := this.ParseForm(config);
	err != nil {
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

func (this *MainController) QiniuUpload() {
	sess := this.GetSession("Username")
	if sess == nil{
		this.Redirect("/login",302)
	}
	this.Data["Username"] = sess
	data := models.SiteConfigMap()
	data["Host"] = "api.qiniu.com"
	data["Parameter"] = "/v6/domain/list?tbl="+data["Bucket"]
	data["Url"] = "http://"+data["Host"]+data["Parameter"]
	Bucket := tools.GetBucketData(data)
	r,_ := regexp.Compile("\"([^\"]*)\"")
	match := r.FindString(Bucket)
	match = strings.Replace(match,"\"","",-1)
	this.Data["Bucket"] = match
	data["Host"] = "rsf.qbox.me"
	data["Parameter"] = "/list?bucket="+data["Bucket"]+"&limit=1000&prefix="
	data["Url"] = "http://"+data["Host"]+data["Parameter"]
	list := tools.GetBucketData(data)
	list = tools.ConvertToString(list,"GB18030","gbk")
	this.Data["list"] = list
	this.Layout = "layout.html"
	this.TplName = "qiniuUpload.html"
}
