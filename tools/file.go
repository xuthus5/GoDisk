package tools

import (
	"io"
	"os"
)

/*************
自定义文件操作
**************/

/* 文件删除 */
func FileRemove(path string) error {
	return os.Remove(path)
}

/* 创建目录 */
func DirCreate(path string) (bool, error) {
	_, err := os.Stat(path) //检测路径状态
	if err == nil {
		return true, nil //没有错误 表明文件夹路径存在 返回true nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

/* 文件移动 */
func FileMove(source, target string) error {
	srcFile, err := os.Open(source) //打开源文件
	if err != nil {
		return err
	}
	defer srcFile.Close()
	tagFile, err := os.Create(target) //打开目标文件
	if err != nil {
		return err
	}
	defer tagFile.Close()
	_, err = io.Copy(tagFile, srcFile) //文件拷贝
	if err != nil {
		return err //操作失败
	}
	code := FileRemove(source) //源文件删除
	if code != nil {
		return code //操作失败
	} else {
		return nil
	}
}
