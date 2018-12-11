/***********************

	api接口

************************/

package controllers

import (
	"GoDisk/models"
	"GoDisk/tools"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/upyun/go-sdk/upyun"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
	var data *ResultData
	if err != nil {
		data = &ResultData{Error: 1, Title: "失败:", Msg: "添加失败！"}
	} else {
		data = &ResultData{Error: 0, Title: "成功:", Msg: "添加成功！"}
	}
	this.Data["json"] = data
	this.ServeJSON()
}

// 修改分类api  路由 /api/category/update
func (this *ApiController) CategoryUpdate() {
	id := this.GetString("id")
	data := &models.Category{}
	info := &ResultData{}
	if err := this.ParseForm(data); err != nil {
		info = &ResultData{Error: 1, Title: "失败:", Msg: "接收表单数据出错！"}
	} else {
		data.Id = tools.StringToInt(id)
		err := models.UpdateCategory(data)
		if err != nil {
			info = &ResultData{Error: 1, Title: "失败:", Msg: "数据库操作出错！"}
		} else {
			info = &ResultData{Error: 0, Title: "成功:", Msg: "修改成功！"}
		}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 删除分类api  路由 /api/category/delete
func (this *ApiController) CategoryDelete() {
	info := &ResultData{}
	//先判断分类数目 为1时，不允许删除
	count, _ := models.TableNumber("category")
	if count == 1 {
		info = &ResultData{Error: 1, Title: "失败:", Msg: "必须保留一个分类！"}
	} else {
		id, _ := strconv.Atoi(this.GetString("id"))
		err := models.DeleteCategory(id)
		if err != nil {
			info = &ResultData{Error: 1, Title: "失败:", Msg: "数据库操作出错！"}
		} else {
			info = &ResultData{Error: 0, Title: "成功:", Msg: "删除成功！"}
		}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 分类列表 路由 /api/category/list
func (this *ApiController) CategoryList() {
	this.Data["json"] = &JsonData{Code: 0, Count: 100, Msg: "", Data: models.GetCategoryJson()}
	this.ServeJSON()
}

//文件部分的API 在这里 定义了后台的所有文件操作请求处理方法
// 这是本地上传的

// 文件上传api 路由 /api/file/upload 返回一个包含文件存储信息的json数据
func (this *ApiController) FileUpload() {
	info := &ResultData{Error: 1, Title: "失败:", Msg: "上传失败！"}
	f, h, err := this.GetFile("file")
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer f.Close()
	//获取当前年月日
	year, month, _ := tools.EnumerateDate()
	savePath := "file/" + year + "/" + month + "/"
	//创建存储目录
	tools.DirCreate(savePath)
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
			info = &ResultData{Error: 1, Title: "结果:", Msg: "上传失败！"}
		} else {
			info = &ResultData{Error: 0, Title: "结果:", Msg: "上传成功！", Data: models.FileInfo(id)}
		}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

//文件列表api 路由 /api/file/list
func (this *ApiController) FileList() {
	this.Data["json"] = &JsonData{Code: 0, Count: 100, Msg: "", Data: models.GetFileJson()}
	this.ServeJSON()
}

// 文件删除 路由 /api/file/delete
func (this *ApiController) FileDelete() {
	info := &ResultData{}
	id, _ := strconv.Atoi(this.GetString("id"))
	//数据库文件删除
	filePath, err := models.FileDelete(id)
	if err != nil {
		info = &ResultData{Error: 1, Title: "失败:", Msg: "数据库操作出错！"}
	} else {
		info = &ResultData{Error: 0, Title: "成功:", Msg: "删除成功！"}
	}
	//本地文件删除
	tools.FileRemove(filePath)
	this.Data["json"] = info
	this.ServeJSON()
}

// 七牛云存储的API操作
// 目前功能比较简单 实现了预览 删除操作

// 七牛云文件上传接口 路由 /api/upload/qiniu
func (this *ApiController) QiniuUpload() {
	f, h, err := this.GetFile("attachment")
	if err != nil {
		log.Fatal("error: ", err)
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
	this.SaveToFile("attachment", filePath)                                   //保存文件到本地
	res := tools.QiniuApi(filePath, saveName, models.RetGroupConfig("Qiniu")) //上传到七牛云
	var data *ResultData
	if res == true {
		data = &ResultData{Error: 1, Title: "结果:", Msg: "上传成功！"}
	} else {
		data = &ResultData{Error: 0, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
	}
	os.Remove(filePath) //移除本地文件
	this.Data["json"] = data
	this.ServeJSON()
}

// 七牛云文件列表接口 路由 /api/file/qiniu/list
func (this *ApiController) QiniuList() {
	data := models.RetGroupConfig("Qiniu")
	data["Host"] = "api.qiniu.com"
	data["Parameter"] = "/v6/domain/list?tbl=" + data["Bucket"]
	data["Url"] = "http://" + data["Host"] + data["Parameter"]
	Bucket := tools.GetBucketData(data)
	r, _ := regexp.Compile("\"([^\"]*)\"")
	match := r.FindString(string(Bucket))
	match = strings.Replace(match, "\"", "", -1)
	data["Host"] = "rsf.qbox.me"
	data["Parameter"] = "/list?bucket=" + data["Bucket"]
	data["Url"] = "http://" + data["Host"] + data["Parameter"]
	body := tools.GetBucketData(data)
	var res Response
	json.Unmarshal([]byte(body), &res)
	this.Data["json"] = JsonData{Msg: match, Data: res.Items}
	this.ServeJSON()
}

// 七牛云文件删除 路由 /api/file/qiniu/delete
func (this *ApiController) QiniuDeleteFile() {
	code := this.GetString("code")
	code = base64.StdEncoding.EncodeToString([]byte(code))
	code = strings.Replace(code, "/", "_", -1)
	code = strings.Replace(code, "+", "-", -1)
	data := models.RetGroupConfig("Qiniu")
	data["Host"] = "rs.qiniu.com"
	data["Parameter"] = "/delete/" + code
	data["Url"] = "http://" + data["Host"] + data["Parameter"]
	var res ResponseError
	json.Unmarshal([]byte(tools.DeleteFile(data)), &res)
	this.Data["json"] = JsonData{Data: res.Error}
	this.ServeJSON()
}

// 又拍云系列API

//又拍云文件列表 路由 /api/file/upyun/list
func (this *ApiController) UpyunList() {
	data := models.RetGroupConfig("Upyun")
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   data["Bucket"],
		Operator: data["Operator"],
		Password: data["Password"],
	})
	list := tools.AllUpyunList(up, "/")
	this.Data["json"] = JsonData{Data: list}
	this.ServeJSON()
}

// 又拍云上传 路由 /api/upload/upyun
func (this *ApiController) UpyunUpload() {
	data := models.RetGroupConfig("Upyun")
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   data["Bucket"],
		Operator: data["Operator"],
		Password: data["Password"],
	})

	f, h, err := this.GetFile("attachment")
	if err != nil {
		log.Fatal("error: ", err)
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
	this.SaveToFile("attachment", filePath) //保存文件到本地
	//上传又拍云
	err = up.Put(&upyun.PutObjectConfig{
		Path:      "/" + saveName,
		LocalPath: filePath,
	})
	var info *ResultData
	if err == nil {
		info = &ResultData{Error: 0, Title: "结果:", Msg: "上传成功！"}
	} else {
		info = &ResultData{Error: 1, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
	}
	os.Remove(filePath) //移除本地文件
	this.Data["json"] = info
	this.ServeJSON()
}

//又拍云删除 路由 /api/file/upyun/delete
func (this *ApiController) UpyunDeleteFile() {
	data := models.RetGroupConfig("Upyun")
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   data["Bucket"],
		Operator: data["Operator"],
		Password: data["Password"],
	})
	err := up.Delete(&upyun.DeleteObjectConfig{
		Path:  this.GetString("path"),
		Async: true,
	})
	info := JsonData{}
	if err != nil {
		info = JsonData{Code: 1}
	} else {
		info = JsonData{Code: 0}
	}
	this.Data["json"] = info
	this.ServeJSON()
}

// 腾讯云存储API //

// 腾讯云文件列表 路由 /api/file/cos/list
func (this *ApiController) CosList() {
	data := models.RetGroupConfig("COS")
	info := &ResultData{Error: 0}
	u, _ := url.Parse("http://" + data["Bucket"] + "-" + data["APPID"] + ".cos." + data["Region"] + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  data["SecretID"],
			SecretKey: data["SecretKey"],
		},
	})

	opt := &cos.BucketGetOptions{
		MaxKeys: 1000,
	}
	v, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		info = &ResultData{Error: 1}
	} else {
		info = &ResultData{Data: v.Contents}
	}

	this.Data["json"] = info
	this.ServeJSON()
}

// 腾讯云文件上传 路由 /api/upload/cos
func (this *ApiController) CosUpload() {
	data := models.RetGroupConfig("COS")
	info := &ResultData{Error: 0}
	//上传的文件示例
	f, h, err := this.GetFile("attachment")
	if err != nil {
		info = &ResultData{Error: 1}
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
	this.SaveToFile("attachment", filePath) //保存文件到本地

	// 创建COS Client实例。
	u, _ := url.Parse("http://" + data["Bucket"] + "-" + data["APPID"] + ".cos." + data["Region"] + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  data["SecretID"],
			SecretKey: data["SecretKey"],
		},
	})

	// 上传本地文件。
	//对象键（Key）是对象在存储桶中的唯一标识。
	//例如，在对象的访问域名 ` bucket1-1250000000.cos.ap-guangzhou.myqcloud.com/test/objectPut.go ` 中，对象键为 test/objectPut.go
	stream, err := os.Open(filePath)
	if err != nil {
		info = &ResultData{Error: 1}
	}
	_, err = c.Object.Put(context.Background(), saveName, stream, nil)

	if err != nil {
		info = &ResultData{Error: 1, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
	}
	os.Remove(filePath) //移除本地文件
	this.Data["json"] = info
	this.ServeJSON()
}

//腾讯云文件删除 路由 /api/file/cos/delete
func (this *ApiController) CosDeleteFile() {
	data := models.RetGroupConfig("COS")
	info := &ResultData{Error: 0}
	u, _ := url.Parse("http://" + data["Bucket"] + "-" + data["APPID"] + ".cos." + data["Region"] + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  data["SecretID"],
			SecretKey: data["SecretKey"],
		},
	})

	objectName := this.GetString("key")
	_, err := c.Object.Delete(context.Background(), objectName)
	if err != nil {
		info = &ResultData{Error: 1}
	}

	this.Data["json"] = info
	this.ServeJSON()

}

// 阿里云存储API //

// 阿里云文件列表 路由 /api/file/oss/list
func (this *ApiController) OssList() {
	data := models.RetGroupConfig("OSS")
	// 创建OSSClient实例。
	client, err := oss.New(data["Endpoint"], data["Accesskey"], data["Secretkey"])
	if err != nil {
		log.Fatal(err)
	}

	// 获取存储空间。
	bucketName := data["Bucket"]
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Fatal(err)
	}

	// 列举所有文件。
	marker := ""
	lsRes, err := bucket.ListObjects(oss.Marker(marker))
	this.Data["json"] = JsonData{Data: lsRes.Objects}
	this.ServeJSON()
}

// 阿里云文件上传 路由 /api/upload/oss
func (this *ApiController) OssUpload() {
	data := models.RetGroupConfig("OSS")
	info := &ResultData{}
	//上传的文件示例
	f, h, err := this.GetFile("attachment")
	if err != nil {
		log.Fatal("error: ", err)
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
	this.SaveToFile("attachment", filePath) //保存文件到本地

	// 创建OSSClient实例。
	client, err := oss.New(data["Endpoint"], data["Accesskey"], data["Secretkey"])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 获取存储空间。
	bucket, err := client.Bucket(data["Bucket"])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 上传本地文件。
	err = bucket.PutObjectFromFile(saveName, filePath)
	if err != nil {
		info = &ResultData{Error: 1, Title: "结果:", Msg: "认证失败！请确保配置信息正确"}
	} else {
		info = &ResultData{Error: 0, Title: "结果:", Msg: "上传成功！"}
	}
	os.Remove(filePath) //移除本地文件
	this.Data["json"] = info
	this.ServeJSON()
}

//阿里云文件删除
func (this *ApiController) OssDeleteFile() {
	data := models.RetGroupConfig("OSS")
	info := &ResultData{Error: 0}
	// 创建OSSClient实例。
	client, err := oss.New(data["Endpoint"], data["Accesskey"], data["Secretkey"])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	bucketName := data["Bucket"]
	objectName := this.GetString("key")

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		info = &ResultData{Error: 1}
	}

	// 删除单个文件。
	err = bucket.DeleteObject(objectName)
	if err != nil {
		info = &ResultData{Error: 1}
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
	info := &ResultData{Error: 0, Title: "成功:", Msg: "信息重置成功！"}
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
		Addition = "Qiniu"
	} else if submit == "upyunInfo" {
		data = &models.UpyunConfigOption{}
		Addition = "Upyun"
	} else if submit == "ossInfo" {
		data = &models.OssConfigOption{}
		Addition = "OSS"
	} else if submit == "cosInfo" {
		data = &models.CosConfigOption{}
		Addition = "COS"
	}
	if err := this.ParseForm(data); err != nil {
		info = &ResultData{Error: 1, Title: "失败:", Msg: "接收表单数据出错！"}
	} else {
		t := reflect.TypeOf(data).Elem()  //类型
		v := reflect.ValueOf(data).Elem() //值
		for i := 0; i < t.NumField(); i++ {
			config := &models.Config{Option: t.Field(i).Name, Value: v.Field(i).String(), Addition: Addition}
			err := models.SiteConfig(config)
			if err != nil {
				info = &ResultData{Error: 1, Title: "失败:", Msg: "出现未知错误！"}
				break
			}
		}
	}
	this.Data["json"] = info
	this.ServeJSON()
}
