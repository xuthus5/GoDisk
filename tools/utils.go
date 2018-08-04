package tools

import (
	"time"
	"crypto/md5"
	"encoding/hex"
	"os"
	"io"
	"log"
	)

/*
时间转字符串
*/
func TimeToString() string{
	//timeLayout := "2006-01-02 15:04:05"                  //时间模板-精确
	timeLayout := "2006-01-02"
	nowTime := time.Now().Unix()                         //当前时间
	dateTime := time.Unix(nowTime, 0).Format(timeLayout) //转换当前时间戳为时间模板格式
	return dateTime	//返回时间字符串
}
/*
字符串转md5
*/
func StringToMd5(str string) string{
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

/*
文件移动
 */
func FileMove(source,target string) bool {
	srcFile, err := os.Open(source)	//打开源文件
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	tagFile, err := os.Create(target)	//打开目标文件
	if err != nil {
		log.Fatal(err)
	}
	defer tagFile.Close()
	_,err = io.Copy(tagFile, srcFile)	//文件拷贝
	if err != nil{
		return false	//操作失败
	}
	code := FileRemove(source)	//源文件删除
	if code == false {
		return false	//操作失败
	}else{
		return true
	}
}

/*
文件删除
 */
func FileRemove(path string) bool {
	err := os.Remove(path)
	if err != nil{
		log.Fatal(err)
		return false	//操作失败
	}else{
		return true	//操作成功
	}
}

/*
创建目录
 */
func DirCreate(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return true
		} else {
			return false
		}
	}
	return true
}