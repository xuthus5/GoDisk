package controllers

import (
	"GoDisk/models"
	"GoDisk/tools"
	"encoding/json"
	"github.com/astaxie/beego"
	"regexp"
	"strings"
	"log"
)

type QiNiuController struct {
	beego.Controller
}

//七牛云资源列表
type List struct {
	Key 	string `json:"key"`
	Hash  	string `json:"hash"`
	Fsize	int64 `json:"fsize"`
	MimeType	string `json:"mimeType"`
	PutTime	int64 `json:"putTime"`
	Type	int64 `json:"type"`
	Status	int64 `json:"status"`
}
type Response struct {
	Items	[]List	`json:"items"`
}

// 七牛云上传页面 包含列表 后续增加文件操作
func (this *QiNiuController) Index() {
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
	match := r.FindString(string(Bucket))
	match = strings.Replace(match,"\"","",-1)
	this.Data["Bucket"] = match
	data["Host"] = "rsf.qbox.me"
	data["Parameter"] = "/list?bucket="+data["Bucket"]
	data["Url"] = "http://"+data["Host"]+data["Parameter"]
	body := tools.GetBucketData(data)
	var res Response
	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		log.Printf("err was %v", err)
	}
	this.Data["list"] = res.Items
	this.Layout = "layout.html"
	this.TplName = "qiniu.html"
}