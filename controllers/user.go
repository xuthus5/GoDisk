/***********************

	用户管理

************************/

package controllers

import (
	"GoDisk/models"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Login() {
	Username := this.GetString("username")
	Password := this.GetString("password")
	if Username == "" || Password == "" {
		this.TplName = "login.html"
	} else {
		user := &models.Config{Option: "Author", Value: Username}
		pass := &models.Config{Option: "Password", Value: Password}
		err := models.Login(user, pass)
		var data *ResultData
		if err == nil {
			this.SetSession("master", Username)
			data = &ResultData{Error: 0, Title: "你好啊", Msg: "欢迎回来！"}
		} else {
			data = &ResultData{Error: 1, Title: "不好啦", Msg: "账号或密码输入有误！"}
		}
		this.Data["json"] = data
		this.ServeJSON()
	}
}

func (this *UserController) Logout() {
	sess := this.GetSession("master")
	if sess != nil {
		//删除指定的session
		this.DelSession("Username")
		//销毁全部的session
		//this.DestroySession()
		this.Redirect("/login", 302)
	}
}
