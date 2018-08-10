package tools

import (
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"context"
	"net/http"
		"os"
	"io/ioutil"
		"log"
)


//上传接口
func QiniuApi(filePath,fileName string,config map[string]string) bool {
	var (
		Accesskey = config["Accesskey"]
		Secretkey = config["Secretkey"]
		Bucket = config["Bucket"]
		Zone = config["Zone"]
	)
	localFile := filePath
	key := fileName

	putPolicy := storage.PutPolicy{
		Scope:Bucket,
	}

	mac := qbox.NewMac(Accesskey, Secretkey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	switch Zone {
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
			"x:name": "GoD upload",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return false
	}
	return true
}

//管理凭证生成
func GeneratingVoucher(data map[string]string) string{
	var (
		Accesskey = data["Accesskey"]
		Secretkey = data["Secretkey"]
	)

	signingStr := data["Parameter"]+"\n"
	signByte := []byte(signingStr)
	mac := qbox.NewMac(Accesskey,Secretkey)
	sign := mac.Sign(signByte)
	return sign
}

//获取Bucket的数据
func GetBucketData(data map[string]string) []byte{
	client := &http.Client{}
	request, err := http.NewRequest("GET", data["Url"], nil) //建立一个请求
	if err != nil {
		log.Println("Fatal error ", err.Error())
		os.Exit(0)
	}

	//获取凭证
	sign := GeneratingVoucher(data)

	//Add 头协议
	request.Header.Add("Host",data["Host"])
	request.Header.Add("User-Agent","Go-http-client/1.1")
	request.Header.Add("Authorization", "QBox "+sign)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request) //提交
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	return body //网页源码
}