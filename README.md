# GoD

> 一个基于beego构建的web存储应用，帮你快速部署存储服务

##  1. 如何使用这个服务

&ensp; 1. 首先，确认已经安装好golang语言环境

&ensp; 2. 执行一下命令，安装一些依赖包与框架(建议使用[gopm](https://gopm.io/)进行包管理)

```bash
go get github.com/astaxie/beego
go get github.com/mattn/go-sqlite3
go get github.com/jmoiron/sqlx
go get github.com/qiniu/api.v7

# 官方协助快速开发工具 bee
go get github.com/beego/bee
```

假若网络原因，无法使用以上命令，请单独下载[master.zip](https://gitee.com/xuthus5/GoDisk/attach_files)资源包，直接解压放到 ```$GOPATH/``` 目录下即可，里面包含有所有的依赖

&ensp; 3. 启动项目

&ensp;&ensp; 进入项目 执行 ```go run main.go``` 启动项目，通过bee工具，在项目下执行 ```bee run```

&ensp; 4. 访问 (默认端口8080,更改端口，请修改 conf/app.conf 》httpport选项)

&ensp;&ensp; 访问： http://ip:8080

# 演示地址

+ [后台演示](http://xblogs.cn:8080/login)
+ 账号密码 admin/admin(请勿修改账户信息)


# 演示截图

![本地上传](http://dl.xuthus.cc/godisk-local.png)
![七牛云上传](http://dl.xuthus.cc/godisk-qiniu.png)
![网站配置](http://dl.xuthus.cc/godisk-set.png)