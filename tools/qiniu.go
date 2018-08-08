package tools

import (
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"context"
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

//获取Bucket空间名称
