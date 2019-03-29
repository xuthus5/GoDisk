/************************

	网站配置

*************************/

package models

import (
	"log"
)

// 网站全局配置 用于更新一条网站配置信息
func SiteConfig(data *Config) error {
	raw := new(Config)
	raw = &Config{Option: data.Option, Value: data.Value, Addition: data.Addition}
	if _, _, err := dbc.ReadOrCreate(data, "Option","Addition"); err != nil {
		return err
	} else {
		_, err := dbc.Raw("UPDATE config SET value = ? WHERE option = ? AND addition = ? ", raw.Value, raw.Option, raw.Addition).Exec()
		return err
	}
}

//返回网站配置信息为map
func SiteConfigMap() map[string]string {
	config := []Config{}
	err := dbx.Select(&config, "select * from config")
	if err != nil {
		log.Fatal(err.Error())
	}
	var data = make(map[string]string)
	for _, v := range config {
		data[v.Option] = v.Value
	}
	return data
}

//添加配置信息
func AddConfig(info *Config) error {
	_, err := dbc.Insert(info)
	return err
}

// 用于获取网站配置信息
func GetOneConfig(info, addition string) string {
	data := []Config{}
	err := dbx.Select(&data, "select * from config where Option='"+info+"' and addition='"+addition+"'")
	if err != nil {
		panic(err.Error())
	}
	return data[0].Value
}

// 返回一个组配置信息
func RetGroupConfig(groupName string) map[string]string {
	config := []Config{}
	err := dbx.Select(&config, "select * from config where addition=?", groupName)
	if err != nil {
		log.Fatal(err.Error())
	}
	var data = make(map[string]string)
	for _, v := range config {
		data[v.Option] = v.Value
	}
	return data
}