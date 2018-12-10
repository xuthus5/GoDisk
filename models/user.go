/************************

	用户管理

*************************/

package models

import "qiniupkg.com/x/log.v7"

// 用户登陆校验
func Login(username, password *Config) error {
	user := dbc.Read(username, "Option", "Value", "Addition")
	pass := dbc.Read(password, "Option", "Value", "Addition")
	if user != nil {
		log.Println(user)
		return user
	}
	if pass != nil {
		log.Println(pass)
		return pass
	}
	return nil
}
