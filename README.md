# GoDisk

> GoDisk是一个基于beego框架构建的web存储应用，能帮你快速部署存储服务。目前已集成七牛云,又拍云,阿里云OSS,腾讯云COS等对象存储

## 如何使用这个服务

### 1. 安装好golang环境

### 2. 执行以下命令

```bash
# 确保已安装如下包
go get github.com/astaxie/beego
go get github.com/mattn/go-sqlite3
go get github.com/jmoiron/sqlx
# 七牛云对象存储API包
go get github.com/qiniu/api.v7
# 又拍云对象存储API包
go get github.com/upyun/go-sdk/upyun
# 腾讯云对象存储API包
go get -u github.com/tencentyun/cos-go-sdk-v5
# 阿里云对象存储API包
go get -u github.com/aliyun/aliyun-oss-go-sdk/oss

# 官方协助快速开发工具 bee[非必须]
go get github.com/beego/bee
```

### 3. 启动项目

```
#直接运行
go run main.go
#通过bee工具快速运行
bee run
```

### 4. 访问
 
默认端口80,更改端口，请修改 conf/app.conf -> runmode 选项;prod:80  dev:8080

# 演示地址

+ [演示](http://xblogs.cn:8080)
+ 账号密码 admin/admin(请勿修改账户信息)