package tools

import (
	"github.com/upyun/go-sdk/upyun"
	"time"
)

type UpyunList struct {
	Name string
	Size int64
	Time time.Time
	Path string
}

func AllUpyunList(up *upyun.UpYun, path string) []UpyunList {
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
				list = append(list, AllUpyunList(up, path+obj.Name+"/")...)
			}
		}
	}
	return list
}
