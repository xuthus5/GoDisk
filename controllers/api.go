/***********************

	api接口

************************/

package controllers

import (
	"GoDisk/models"
	"GoDisk/tools"
	"encoding/json"
	"github.com/astaxie/beego"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"
)

type ApiController struct {
	beego.Controller
}

// 分类管理 分类管理用于对上传的文件进行分类 目前尚未对分类进行使用 后续的升级过程中 将会加入分类功能
// 使文件管理变得更加的完美

// 添加分类api  路由 /api/category/add
func (this *ApiController) CategoryAdd() {
	name := this.GetString("name")
	key := this.GetString("key")
	description := this.GetString("description")
	info := &models.Category{Name: name, Key: key, Description: description}
	err := models.AddCategory(info)
	var data *Result
	if err != nil {
		data = &Result{Error: 1, Title: "失败:", Msg: "添加失败！"}
	} else {
		data = &Result{Error: 0, Title: "成功:", Msg: "添加成功！"}
	}
	this.Data["json"] = data
	this.ServeJSON()
}

// 修改分类api  路由 /api/category/update
func (this *ApiController) CategoryUpdate() {
	id := this.GetString("id")
	data := &models.Category{}
	info := &Result{}
	if err := this.ParseForm(data); err != nil {
		info = &Result{Error: 1, Title: "失败:", Msg: "接收表单数据出错！"}
	} else {
		data.Id = tools.StringToInt(id)
		err := models.UpdateCategory(data)
		if err != nil {
			info = &Result{Error: 1, Title: "失败:", Msg: "数据库操作出错！"}
		} else {
			info = &Result{Error: 0, Title: "成功:", Msg: "修改成功！"}
		}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 删除分类api  路由 /api/category/delete
func (this *ApiController) CategoryDelete() {
	info := &Result{}
	//先判断分类数目 为1时，不允许删除
	count, _ := models.TableNumber("category")
	if count == 1 {
		info = &Result{Error: 1, Title: "失败:", Msg: "必须保留一个分类！"}
	} else {
		id, _ := strconv.Atoi(this.GetString("id"))
		err := models.DeleteCategory(id)
		if err != nil {
			info = &Result{Error: 1, Title: "失败:", Msg: "数据库操作出错！"}
		} else {
			info = &Result{Error: 0, Title: "成功:", Msg: "删除成功！"}
		}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 分类列表 路由 /api/category/list
func (this *ApiController) CategoryList() {
	this.Data["json"] = &Result{Error: 0, Count: 100, Msg: "", Data: models.GetCategoryJson()}
	this.ServeJSON()
}

//文件部分的API 在这里 定义了后台的所有文件操作请求处理方法
// 这是本地上传的

// 文件上传api 路由 /api/file/upload 返回一个包含文件存储信息的json数据
func (this *ApiController) FileUpload() {
	info := &Result{Error: 1, Title: "失败:", Msg: "上传失败！"}
	f, h, err := this.GetFile("file")
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer f.Close()
	//获取当前年月日
	year, month, _ := tools.EnumerateDate()
	savePath := "file/" + year + "/" + month + "/"
	//创建存储目录
	_, _ = tools.DirCreate(savePath)
	//重命名文件名称
	tempFileName := tools.StringToMd5(h.Filename, 5)
	suffix := tools.GetFileSuffix(h.Filename)
	saveName := tempFileName + suffix
	// 保存位置
	err = this.SaveToFile("file", savePath+saveName)
	//写入数据库
	if err == nil {
		//写入数据库
		data := &models.Attachment{Name: saveName, Path: savePath + saveName, Created: tools.Int64ToString(time.Now().Unix())}
		id, code := models.FileSave(data)
		if code != nil {
			info = &Result{Error: 1, Title: "结果:", Msg: "上传失败！"}
		} else {
			info = &Result{Error: 0, Title: "结果:", Msg: "上传成功！", Data: models.FileInfo(id)}
		}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

//文件列表api 路由 /api/file/list
func (this *ApiController) FileList() {
	this.Data["json"] = &Result{Error: 0, Count: 100, Msg: "", Data: models.GetFileJson()}
	this.ServeJSON()
}

// 文件删除 路由 /api/file/delete
func (this *ApiController) FileDelete() {
	info := &Result{}
	id, _ := strconv.Atoi(this.GetString("id"))
	//数据库文件删除
	filePath, err := models.FileDelete(id)
	if err != nil {
		info = &Result{Error: 1, Title: "失败:", Msg: "数据库操作出错！"}
	} else {
		info = &Result{Error: 0, Title: "成功:", Msg: "删除成功！"}
	}
	//本地文件删除
	_ = tools.FileRemove(filePath)
	this.Data["json"] = info
	this.ServeJSON()
}

// 七牛云存储的API操作
// 目前功能比较简单 实现了预览 删除操作

// 七牛云文件上传接口 路由 /api/upload/qiniu
func (this *ApiController) QiniuUpload() {
	//文件上传
	f, h, err := this.GetFile("attachment")
	defer f.Close()
	if err != nil {
		this.Data["json"] = &Result{Error: 0, Title: "结果:", Msg: err.Error()}
		this.ServeJSON()
	} else {
		fileName := this.GetString("customName") //自定义文件名
		saveName := ""                           //文件存储名
		if fileName == "" {
			saveName = h.Filename
		} else {
			fileSuffix := path.Ext(h.Filename) //得到文件后缀
			saveName = fileName + fileSuffix
		}
		filePath := "file/" + saveName
		_ = this.SaveToFile("attachment", filePath) //保存文件到本地
		//文件转储成功 上传远端
		config := models.RetGroupConfig("qn")
		factory := tools.EntityFactory{}
		err := factory.Create("qn", tools.Qiniu{Accesskey: config["QnAk"], Secretkey: config["QnSk"], Zone: config["QnZone"], Bucket: config["QnBucket"]}).Upload(filePath, saveName)
		var data *Result
		_ = os.Remove(filePath) //移除本地文件
		if err == nil {
			data = &Result{Error: 1, Title: "结果:", Msg: "上传成功！"}
		} else {
			data = &Result{Error: 0, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
		}
		this.Data["json"] = data
		this.ServeJSON()
	}
}

// 七牛云文件列表接口 路由 /api/file/qiniu/list
func (this *ApiController) QiniuList() {
	config := models.RetGroupConfig("qn")
	qn := tools.Qiniu{
		Accesskey: config["QnAk"],
		Secretkey: config["QnSk"],
		Zone:      config["QnZone"],
		Bucket:    config["QnBucket"],
	}
	var res Response
	var info Result
	err, body, bucket := qn.List()
	if err != nil {
		info = Result{Error: 0, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
	} else {
		_ = json.Unmarshal([]byte(body), &res)
		info = Result{Error: 1, Title: "结果:", Msg: bucket, Data: res.Items}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 七牛云文件删除 路由 /api/file/qiniu/delete
func (this *ApiController) QiniuDeleteFile() {
	config := models.RetGroupConfig("qn")
	factory := tools.EntityFactory{}
	qn := factory.Create("qn", tools.Qiniu{Accesskey: config["QnAk"], Secretkey: config["QnSk"], Zone: config["QnZone"], Bucket: config["QnBucket"], Host: "rs.qiniu.com"})
	err := qn.Delete(this.GetString("code"))
	info := Result{}
	if err.Error() != "" {
		info = Result{Error: 1, Msg: err.Error()}
	} else {
		info = Result{Error: 0}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 又拍云系列API

//又拍云文件列表 路由 /api/file/upyun/list
func (this *ApiController) UpyunList() {
	data := models.RetGroupConfig("up")
	up := tools.Upyun{Bucket: data["UpBucket"], Operator: data["UpOperator"], Password: data["UpPassword"]}
	list := up.List("/")
	this.Data["json"] = Result{Data: list}
	this.ServeJSON()
}

// 又拍云上传 路由 /api/upload/upyun
func (this *ApiController) UpyunUpload() {
	f, h, err := this.GetFile("attachment")
	if err != nil {
		this.Data["json"] = &Result{Error: 0, Title: "结果:", Msg: err.Error()}
		this.ServeJSON()
	}
	defer f.Close()
	fileName := this.GetString("customName") //自定义文件名
	saveName := ""                           //文件存储名
	if fileName == "" {
		saveName = h.Filename
	} else {
		fileSuffix := path.Ext(h.Filename) //得到文件后缀
		saveName = fileName + fileSuffix
	}
	filePath := "file/" + saveName
	_ = this.SaveToFile("attachment", filePath) //保存文件到本地
	data := models.RetGroupConfig("up")
	factory := tools.EntityFactory{}
	err = factory.Create("up", tools.Upyun{Bucket: data["UpBucket"], Operator: data["UpOperator"], Password: data["UpPassword"]}).Upload("/"+saveName, filePath)
	//上传又拍云
	var info *Result
	if err == nil {
		info = &Result{Error: 0, Title: "结果:", Msg: "上传成功！"}
	} else {
		info = &Result{Error: 1, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
	}
	_ = os.Remove(filePath) //移除本地文件
	this.Data["json"] = info
	this.ServeJSON()
}

//又拍云删除 路由 /api/file/upyun/delete
func (this *ApiController) UpyunDeleteFile() {
	data := models.RetGroupConfig("up")
	factory := tools.EntityFactory{}
	err := factory.Create("up", tools.Upyun{Bucket: data["UpBucket"], Operator: data["UpOperator"], Password: data["UpPassword"]}).Delete(this.GetString("path"))
	info := Result{}
	if err != nil {
		info = Result{Error: 1}
	} else {
		info = Result{Error: 0}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 腾讯云存储API //

// 腾讯云文件列表 路由 /api/file/cos/list
func (this *ApiController) CosList() {
	data := models.RetGroupConfig("cos")
	ten := tools.Cos{Bucket: data["CosBucket"], Appid: data["CosAppid"], Region: data["CosRegion"], Sk: data["CosSk"], Skid: data["CosSkid"]}
	err, body := ten.List()
	info := &Result{Error: 0}
	if err != nil {
		info = &Result{Error: 1}
	} else {
		info = &Result{Data: body.Contents}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 腾讯云文件上传 路由 /api/upload/cos
func (this *ApiController) CosUpload() {
	info := &Result{Error: 0}
	//上传的文件示例
	f, h, err := this.GetFile("attachment")
	if err != nil {
		info = &Result{Error: 1}
		this.Data["json"] = info
		this.ServeJSON()
	}
	defer f.Close()
	fileName := this.GetString("customName") //自定义文件名
	saveName := ""                           //文件存储名
	if fileName == "" {
		saveName = h.Filename
	} else {
		fileSuffix := path.Ext(h.Filename) //得到文件后缀
		saveName = fileName + fileSuffix
	}
	filePath := "file/" + saveName
	_ = this.SaveToFile("attachment", filePath) //保存文件到本地
	data := models.RetGroupConfig("cos")
	factory := tools.EntityFactory{}
	err = factory.Create("cos", tools.Cos{Bucket: data["CosBucket"], Appid: data["CosAppid"], Region: data["CosRegion"], Sk: data["CosSk"], Skid: data["CosSkid"]}).Upload(filePath, saveName)
	if err != nil {
		info = &Result{Error: 1, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
	}
	_ = os.Remove(filePath) //移除本地文件
	this.Data["json"] = info
	this.ServeJSON()
}

//腾讯云文件删除 路由 /api/file/cos/delete
func (this *ApiController) CosDeleteFile() {
	data := models.RetGroupConfig("cos")
	factory := tools.EntityFactory{}
	err := factory.Create("cos", tools.Cos{Bucket: data["CosBucket"], Appid: data["CosAppid"], Region: data["CosRegion"], Sk: data["CosSk"], Skid: data["CosSkid"]}).Delete(this.GetString("key"))
	info := &Result{Error: 0}
	if err != nil {
		info = &Result{Error: 1}
	}
	this.Data["json"] = info
	this.ServeJSON()

}

// 阿里云存储API //

// 阿里云文件列表 路由 /api/file/oss/list
func (this *ApiController) OssList() {
	info := &Result{Error: 0}
	data := models.RetGroupConfig("oss")
	// 创建OSSClient实例。
	ali := tools.Oss{Bucket: data["OssBucket"], Ak: data["OssAk"], Sk: data["OssSk"], Endpoint: data["OssEndpoint"]}
	list, err := ali.List()
	if err != nil {
		info = &Result{Error: 1}
	} else {
		info = &Result{Data: list.Objects}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 阿里云文件上传 路由 /api/upload/oss
func (this *ApiController) OssUpload() {
	//上传的文件示例
	f, h, err := this.GetFile("attachment")
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer f.Close()
	info := &Result{Error: 0}
	// 创建OSSClient实例。
	fileName := this.GetString("customName") //自定义文件名
	saveName := ""                           //文件存储名
	if fileName == "" {
		saveName = h.Filename
	} else {
		fileSuffix := path.Ext(h.Filename) //得到文件后缀
		saveName = fileName + fileSuffix
	}
	filePath := "file/" + saveName
	_ = this.SaveToFile("attachment", filePath) //保存文件到本地
	// 上传本地文件。
	data := models.RetGroupConfig("oss")
	factory := tools.EntityFactory{}
	err = factory.Create("oss", tools.Oss{Bucket: data["OssBucket"], Ak: data["OssAk"], Sk: data["OssSk"], Endpoint: data["OssEndpoint"]}).Upload(saveName, filePath)
	if err != nil {
		info = &Result{Error: 1, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
	} else {
		info = &Result{Error: 0, Title: "结果:", Msg: "上传成功！"}
	}
	_ = os.Remove(filePath) //移除本地文件
	this.Data["json"] = info
	this.ServeJSON()
}

//阿里云文件删除
func (this *ApiController) OssDeleteFile() {
	info := &Result{Error: 0}
	data := models.RetGroupConfig("oss")
	factory := tools.EntityFactory{}
	err := factory.Create("oss", tools.Oss{Bucket: data["OssBucket"], Ak: data["OssAk"], Sk: data["OssSk"], Endpoint: data["OssEndpoint"]}).Delete(this.GetString("key"))
	if err != nil {
		info = &Result{Error: 1}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

//网站配置页面的处理信息都在这里
// 通过提交判断submit的值来判断是哪一个表单提交的

// 网站设置页面  路由  /api/site/config
func (this *ApiController) SiteConfig() {
	// 判断提交类型 user为用户信息表单  site为网站配置表单
	submit := this.GetString("submit")
	info := &Result{Error: 0, Title: "成功:", Msg: "信息重置成功！"}
	Addition := ""
	var data interface{}
	if submit == "userInfo" {
		data = &models.UserConfigOption{}
		Addition = ""
	} else if submit == "siteInfo" {
		data = &models.SiteConfigOption{}
		Addition = ""
	} else if submit == "niniuInfo" {
		data = &models.QiniuConfigOption{}
		Addition = "qn"
	} else if submit == "upyunInfo" {
		data = &models.UpyunConfigOption{}
		Addition = "up"
	} else if submit == "ossInfo" {
		data = &models.OssConfigOption{}
		Addition = "oss"
	} else if submit == "cosInfo" {
		data = &models.CosConfigOption{}
		Addition = "cos"
	}
	if err := this.ParseForm(data); err != nil {
		info = &Result{Error: 1, Title: "失败:", Msg: "接收表单数据出错！"}
	} else {
		t := reflect.TypeOf(data).Elem()  //类型
		v := reflect.ValueOf(data).Elem() //值
		for i := 0; i < t.NumField(); i++ {
			config := &models.Config{Option: t.Field(i).Name, Value: v.Field(i).String(), Addition: Addition}
			err := models.SiteConfig(config)
			if err != nil {
				info = &Result{Error: 1, Title: "失败:", Msg: "出现未知错误！"}
				break
			}
		}
	}
	this.Data["json"] = info
	this.ServeJSON()
}
