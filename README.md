# GoD

> 一个基于beego构建的web存储应用，帮你快速部署存储服务

##  1. 如何使用这个服务

&ensp; 1. 首先，确认已经安装好golang语言环境

&ensp; 2. 执行一下命令，安装一些依赖包与框架

```bash
go get github.com/astaxie/beego
go get github.com/mattn/go-sqlite3
go get github.com/jmoiron/sqlx
```

&ensp; 3. 启动项目

&ensp;&ensp; 进入项目 执行 go run main.go 即可启动项目

&ensp;&ensp;你也可以使用beego官方提供的 bee工具 来快速启动，你可以通过如下的方式安装 bee 工具，安装完毕后，进入项目 执行 ```bee run``` 即可快速启动

```bash
# 安装bee工具
go get github.com/beego/bee
```
