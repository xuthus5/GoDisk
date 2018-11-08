/************************

	用户管理

*************************/

package models

// 用户登陆校验
func Login(username, password *Config) error {
	user := dbc.Read(username, "Option", "Value")
	pass := dbc.Read(password, "Option", "Value")
	if user != nil {
		return user
	}
	if pass != nil {
		return pass
	}
	return nil
}
