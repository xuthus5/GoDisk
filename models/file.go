/************************

	文件操作

*************************/

package models

import (
	"GoDisk/tools"
)

//添加分类
func AddCategory(info *Category) error {
	_, err := dbc.Insert(info)
	return err
}

// 删除分类
func DeleteCategory(id int) error {
	//删除分类的同时，需要将旗下的文章删除
	data := &Category{Id: id}
	_, err := dbc.Delete(data)
	if err != nil {
		return err
	}
	all := &Attachment{Cid: id}
	_, err = dbc.Delete(all, "Cid")
	if err != nil {
		return err
	}
	return nil
}

//获取后台分类列表
func GetCategoryJson() *[]Category {
	list := []Category{}
	err := dbx.Select(&list, "select * from category order by id desc")
	if err != nil {
		panic(err.Error())
	}
	return &list
}

// 分类修改
func UpdateCategory(data *Category) error {
	_, err := dbc.Update(data)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 获得一个分类信息 主要用于更新分类信息
func GetOneCategoryInfo(id string) *[]Category {
	list := []Category{}
	dbx.Select(&list, "select * from category where id=?", id)
	return &list
}

//获取后台文件列表
func GetFileJson() *[]Attachment {
	list := []Attachment{}
	err := dbx.Select(&list, "select * from attachment order by id desc")
	if err != nil {
		panic(err.Error())
	}
	for i, v := range list {
		list[i].Created = tools.UnixTimeToString(v.Created)
	}
	return &list
}

//文件上传 入数据库
func FileSave(info *Attachment) (int64, error) {
	return dbc.Insert(info)
}

// 返回一个附件信息
func FileInfo(id int64) *[]Attachment {
	data := []Attachment{}
	err := dbx.Select(&data, "select * from attachment where id=?", id)
	if err != nil {
		panic(err.Error())
	}
	return &data
}

// 文件删除
func FileDelete(id int) (string, error) {
	data := &Attachment{Id: id}
	dbc.Read(data)
	_, err := dbc.Delete(data)
	return data.Path, err
}
