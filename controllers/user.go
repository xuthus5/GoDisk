package controllers

import (
	"github.com/astaxie/beego"
	"GoDisk/models"
	"GoDisk/tools"
	)

type UserController struct {
	beego.Controller
}

type ResultData struct {
	Code  int
	Title string
	Msg   string
}

func (this *UserController) Login() {
	Username := this.GetString("username")
	Password := this.GetString("password")
	if Username == "" || Password == "" {
		this.TplName = "login.html"
	} else {
		user := &models.User{Username:Username,Password:tools.StringToMd5(Password)}
		code,msg := models.Login(user)
		var data *ResultData
		if code == 1{
			this.SetSession("Username",Username)
			data = &ResultData{Code:1,Title:"你好啊",Msg:msg}
		}else{
			data = &ResultData{Code:0,Title:"不好啦",Msg:msg}
		}
		this.Data["json"] = data
		this.ServeJSON()
	}
}

func (this *UserController) Logout() {
	sess := this.GetSession("Username")
	if sess != nil {
		//删除指定的session
		this.DelSession("Username")
		//销毁全部的session
		//this.DestroySession()
		this.Redirect("/login",302)
	}
}
