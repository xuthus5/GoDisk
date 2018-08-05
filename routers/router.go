package routers

import (
	"GoDisk/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//页面路由
    beego.Router("/", &controllers.MainController{})
	beego.Router("/admin",&controllers.MainController{},"*:Admin")
    beego.Router("/classify",&controllers.MainController{},"*:Classify")
    beego.Router("/setting",&controllers.MainController{},"*:Setting")
    beego.Router("/filemanager",&controllers.MainController{},"*:FileManager")

    //用户模块
	beego.Router("/login",&controllers.UserController{},"*:Login")
	beego.Router("/logout",&controllers.UserController{},"*:Logout")

    //接口Api
    beego.Router("/api/upload",&controllers.ApiController{},"post:Upload")
	beego.Router("/api/saveFile",&controllers.ApiController{},"post:SaveFile")
}
