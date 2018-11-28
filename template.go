package main

import (
	"GoDisk/models"
)

/* 直接调取数据库配置表-网站基本信息 */
func SiteConfig(info,addition string) string {
	return models.GetOneConfig(info,addition)
}
