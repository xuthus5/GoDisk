/***********************

	数据结构

************************/

package controllers

// 返回的数据结构
type ResultData struct {
	Error int
	Title string
	Msg   string
	Data  interface{} `json:"data"` //数据
}

// 返回json列表 数据格式
type JsonData struct {
	Code  int         `json:"code"`  //错误代码
	Count int         `json:"count"` // 数据数量
	Msg   string      `json:"msg"`   //输出信息
	Data  interface{} `json:"data"`  //数据
}

//上传成功 返回数据
type UploadData struct {
	Code     int64
	Fid      int64
	FileName string
}

//七牛云资源列表
type List struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	MimeType string `json:"mimeType"`
	PutTime  int64  `json:"putTime"`
	Type     int64  `json:"type"`
	Status   int64  `json:"status"`
}
type Response struct {
	Items []List `json:"items"`
}
// 七牛云删除响应
type ResponseError struct {
	Error 	string `json:"error"`
}