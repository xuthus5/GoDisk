package controllers

import (
	"github.com/astaxie/beego"
	"log"
	"GoDisk/tools"
	"GoDisk/models"
	"path"
	"strings"
	"os"
)

type ApiController struct {
	beego.Controller
}

type UploadData struct {
	Code 	int64
	Fid		int64
	FileName 	string
}

func (this *ApiController) Upload() {
	/*
	备忘 文件与表单异步提交
	文件快表单一步提交
	只有表单提交 文件才能写入数据库
	此处只做文件上传
	移动操作另用函数
	 */
	f, h, err := this.GetFile("file")
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer f.Close()
	savePath := "data/temporary/" + h.Filename
	this.SaveToFile("file", savePath) // 保存位置, 没有文件夹要先创建
	this.Data["json"] = &UploadData{Code:1,FileName:h.Filename}
	this.ServeJSON()
}

func (this *ApiController) SaveFile() {
	//文件存储 表单提交
	fileName := this.GetString("name")		//自定义文件名
	fileMark := this.GetString("mark")		//文件分类
	tempName := this.GetString("filename")	//临时文件名
	saveName := ""								//文件存储名
	saveMark := ""								//文件存储分类
	if fileName == ""{
		saveName = tempName
	}else{
		fileSuffix := path.Ext(tempName)		//得到文件后缀
		fileName = strings.Replace(fileName,".","",-1)
		saveName = fileName+fileSuffix
	}
	if fileMark == ""{
		saveMark = "default"
	}else{
		saveMark = fileMark
	}
	savePath := "data/"+saveMark+"/"+saveName
	err := tools.FileMove("data/temporary/"+tempName,savePath)
	var data *ResultData
	if err == true {
		//写入数据库
		info := &models.File{Name:saveName,Path:savePath,Mark:saveMark,Created:tools.TimeToString()}
		code := models.FileSave(info)
		if code == false {
			data = &ResultData{Code:0,Title:"结果:",Msg:"上传失败！"}
		}else{
			data = &ResultData{Code:1,Title:"结果:",Msg:"上传成功！"}
		}
	}
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *ApiController) QiniuUpload() {
	//七牛文件上传处理
	//文件存储 表单提交
	fileName := this.GetString("name")		//自定义文件名
	tempName := this.GetString("filename")	//临时文件名
	saveName := ""								//文件存储名
	if fileName == ""{
		saveName = tempName
	}else{
		fileSuffix := path.Ext(tempName)		//得到文件后缀
		fileName = strings.Replace(fileName,".","",-1)
		saveName = fileName+fileSuffix
	}
	filePath := "data/temporary/"+tempName
	err := tools.QiniuApi(filePath,saveName,models.SiteConfigMap())
	var data *ResultData
	if err == true {
		data = &ResultData{Code:1,Title:"结果:",Msg:"上传成功！但是你可能需要前往七牛云官网查看，预览功能还没有写"}
	}else{
		data = &ResultData{Code:0,Title:"结果:",Msg:"认证失败！请确保配置信息正确"}
	}
	os.Remove(filePath)
	this.Data["json"] = data
	this.ServeJSON()
}