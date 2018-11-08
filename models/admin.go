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
	raw = &Config{Option: data.Option, Value: data.Value}
	if _, _, err := dbc.ReadOrCreate(data, "Option"); err != nil {
		return err
	} else {
		_, err := dbc.Raw("UPDATE config SET value = ? WHERE option = ?", raw.Value, raw.Option).Exec()
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
func GetOneConfig(info string) string {
	data := []Config{}
	err := dbx.Select(&data, "select * from config where Option=?", info)
	if err != nil {
		panic(err.Error())
	}
	if data[0].Option == "Zone" && data[0].Value == "" {
		return "undefined"
	}
	return data[0].Value
}
