# GO BBS 说明文档
<img alt="GO BBS" src="http://go-bbs.com/image/logo.png">


## Go Bbs 是什么?
Go Bbs 是一个基于 Beego 开发的可切换模板的 BBS 社交博客系统。它安装简单便捷，页面简介优美。前端是HTML+JS+CSS，不需要掌握一些前端技术栈也能轻松自定义页面。

## 官网
http://go-bbs.com/

## DEMO 
https://nihongdengxia.com/ <br>
http://clblog.club/

## 安装

#### 方式一 源码自动安装
1. 把 Go Bbs 项目拉到本地 `git clone https://github.com/1920853199/go-blog.git`
2. 进入项目的根目录下执行 `go build -o go-bbs`
3. 创建一个配置文件，app.conf 内容如 app.conf.example，配置好环境变量 BEEGO_CONFIG_PATH = '配置文件路径' 
4. 执行./go-bbs --install 安装数据库
4. 最后执行 ./go-bbs 访问对应端口即可

> 管理后台默认账号密码：admin,123456
