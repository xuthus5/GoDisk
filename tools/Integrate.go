package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/upyun/go-sdk/upyun"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

//实体操作
type Entity interface {
	Upload(string, string) error
	Delete(string) error
}

//实体工厂
type EntityFactory struct {
}

//工厂制造实体
func (this *EntityFactory) Create(flag string, entity interface{}) Entity {
	//断言?
	switch flag {
	case "qn":
		value, ok := entity.(Qiniu)
		if ok {
			return &value
		}
	case "up":
		value, ok := entity.(Upyun)
		if ok {
			return &value
		}
	case "oss":
		value, ok := entity.(Oss)
		if ok {
			return &value
		}
	case "cos":
		value, ok := entity.(Cos)
		if ok {
			return &value
		}
	}
	return nil
}

type Qiniu struct {
	Accesskey string
	Secretkey string
	Bucket    string
	Zone      string
	Url       string
	Host      string
	Parameter string
}

// 七牛云删除响应
type ResponseError struct {
	Errors string `json:"error"`
}

//上传实现
func (this *Qiniu) Upload(filePath, key string) error {
	putPolicy := storage.PutPolicy{
		Scope: this.Bucket,
	}
	mac := qbox.NewMac(this.Accesskey, this.Secretkey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	switch this.Zone {
	case "storage.ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "storage.ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "storage.ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "storage.ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	default:
		cfg.Zone = &storage.ZoneHuanan
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "GoDisk upload",
		},
	}
	return formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, &putExtra)
}

//删除实现
func (this *Qiniu) Delete(code string) error {
	code = base64.StdEncoding.EncodeToString([]byte(code))
	code = strings.Replace(code, "/", "_", -1)
	code = strings.Replace(code, "+", "-", -1)
	this.Parameter = "/delete/" + code
	this.Url = "http://" + this.Host + this.Parameter
	client := &http.Client{}
	request, err := http.NewRequest("POST", this.Url, nil) //建立一个请求
	if err != nil {
		return err
	}
	//获取凭证
	sign := this.GeneratingVoucher()
	//Add 头协议
	request.Header.Add("Host", this.Host)
	request.Header.Add("User-Agent", "Go-http-client/1.1")
	request.Header.Add("Authorization", "QBox "+sign)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request) //提交
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var res ResponseError
	_ = json.Unmarshal(body, &res)
	return errors.New(res.Errors)
}

//列表
func (this *Qiniu) List() (error, []byte, string) {
	this.Host = "api.qiniu.com"
	this.Parameter = "/v6/domain/list?tbl=" + this.Bucket
	this.Url = "http://" + this.Host + this.Parameter
	err, Bucket := this.GetBucketData()
	if err != nil {
		return err, nil, ""
	}
	r, _ := regexp.Compile("\"([^\"]*)\"")
	match := r.FindString(string(Bucket))
	match = strings.Replace(match, "\"", "", -1)
	this.Host = "rsf.qbox.me"
	this.Parameter = "/list?bucket=" + this.Bucket
	this.Url = "http://" + this.Host + this.Parameter
	err, body := this.GetBucketData()
	if err != nil {
		return err, nil, ""
	}
	return nil, body, match
}

//管理凭证生成
func (this *Qiniu) GeneratingVoucher() string {
	signingStr := this.Parameter + "\n"
	signByte := []byte(signingStr)
	mac := qbox.NewMac(this.Accesskey, this.Secretkey)
	sign := mac.Sign(signByte)
	return sign
}

//获取Bucket的数据
func (this *Qiniu) GetBucketData() (error, []byte) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", this.Url, nil) //建立一个请求
	if err != nil {
		return err, nil
	}
	//获取凭证
	sign := this.GeneratingVoucher()
	//Add 头协议
	request.Header.Add("Host", this.Host)
	request.Header.Add("User-Agent", "Go-http-client/1.1")
	request.Header.Add("Authorization", "QBox "+sign)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request) //提交
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err, nil
	}
	return nil, body //网页源码
}

type Upyun struct {
	Bucket   string //Bucket
	Operator string //Operator
	Password string //Password
	Domain   string //Domain
}
type UpyunList struct {
	Name string
	Size int64
	Time time.Time
	Path string
}

func (this *Upyun) Upload(remote, local string) error {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   this.Bucket,
		Operator: this.Operator,
		Password: this.Password,
	})
	return up.Put(&upyun.PutObjectConfig{
		Path:      remote,
		LocalPath: local,
	})
}
func (this *Upyun) Delete(remote string) error {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   this.Bucket,
		Operator: this.Operator,
		Password: this.Password,
	})
	return up.Delete(&upyun.DeleteObjectConfig{
		Path:  remote,
		Async: false,
	})
}
func (this *Upyun) List(path string) []UpyunList {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   this.Bucket,
		Operator: this.Operator,
		Password: this.Password,
	})
	list := []UpyunList{} //需要返回的文件列表
	objsChan := make(chan *upyun.FileInfo, 10)
	err := up.List(&upyun.GetObjectsConfig{
		Path:        path,
		ObjectsChan: objsChan,
	})
	if err == nil {
		for obj := range objsChan {
			//判断是否为目录
			if obj.IsDir == false {
				list = append(list, UpyunList{Name: obj.Name, Size: obj.Size, Time: obj.Time, Path: path + obj.Name})
			} else {
				list = append(list, this.List(path+obj.Name+"/")...)
			}
		}
	}
	return list
}

// 阿里云
type Oss struct {
	Bucket   string //Bucket
	Ak       string //Accesskey
	Sk       string //Secretkey
	Endpoint string //地域节点
}

func (this *Oss) Upload(key, filePath string) error {
	client, err := oss.New(this.Endpoint, this.Ak, this.Sk)
	if err != nil {
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(this.Bucket)
	if err != nil {
		return err
	}
	return bucket.PutObjectFromFile(key, filePath)
}
func (this *Oss) Delete(key string) error {
	client, err := oss.New(this.Endpoint, this.Ak, this.Sk)
	if err != nil {
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(this.Bucket)
	if err != nil {
		return err
	}
	return bucket.DeleteObject(key)
}
func (this *Oss) List() (oss.ListObjectsResult, error) {
	client, err := oss.New(this.Endpoint, this.Ak, this.Sk)
	if err != nil {
		return oss.ListObjectsResult{}, err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(this.Bucket)
	if err != nil {
		return oss.ListObjectsResult{}, err
	}
	// 列举所有文件。
	marker := ""
	return bucket.ListObjects(oss.Marker(marker))
}

//腾讯云
type Cos struct {
	Bucket string //Bucket
	Appid  string //APPID
	Region string //Region
	Skid   string //SecretID
	Sk     string //SecretKey
}

func (this *Cos) Upload(filePath, saveName string) error {
	u, _ := url.Parse("http://" + this.Bucket + "-" + this.Appid + ".cos." + this.Region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  this.Skid,
			SecretKey: this.Sk,
		},
	})
	stream, err := os.Open(filePath)
	if err != nil {
		return err
	}
	_, err = c.Object.Put(context.Background(), saveName, stream, nil)
	return err
}
func (this *Cos) Delete(objectName string) error {
	u, _ := url.Parse("http://" + this.Bucket + "-" + this.Appid + ".cos." + this.Region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  this.Skid,
			SecretKey: this.Sk,
		},
	})
	_, err := c.Object.Delete(context.Background(), objectName)
	return err
}
func (this *Cos) List() (error, cos.BucketGetResult) {
	u, _ := url.Parse("http://" + this.Bucket + "-" + this.Appid + ".cos." + this.Region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  this.Skid,
			SecretKey: this.Sk,
		},
	})
	opt := &cos.BucketGetOptions{
		MaxKeys: 1000,
	}
	v, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		return err, cos.BucketGetResult{}
	} else {
		return nil, cos.BucketGetResult{Contents: v.Contents}
	}
}
