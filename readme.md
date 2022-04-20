# GO BBS 说明文档
<img alt="GO BBS" src="./logo.png">

## Go Bbs 是什么?
Go Bbs 是一个基于 Beego 开发的可切换模板的 BBS 社交博客系统。它安装简单便捷，页面简介优美。前端是HTML+JS+CSS，不需要掌握一些前端技术栈也能轻松自定义页面。

## 官网
http://go-bbs.com/

## DEMO 
https://nihongdengxia.com/
http://clblog.club/

## 安装

#### 方式一 源码自动安装 （推荐）
1. 把 Go Bbs 项目拉到本地 `git clone https://github.com/1920853199/go-blog.git`
2. 进入项目的根目录下执行 `go build -o go-bbs`
3. 创建一个配置文件，app.conf 内容如 app.conf.example，配置好环境变量 BEEGO_CONFIG_PATH = '配置文件路径' 
4. 执行./go-bbs --install 安装数据库
4. 最后执行 ./go-bbs 访问对应端口即可


#### 方式二 Docker 安装

1. 先安装`docker`以及`docker-compose`
2. 把根目录下的`docker-compose.yml`赋值到你需要运行的`Go Blog`项目的目录下，执行`docker-compose up -d`.（会报找不到数据库的错误，忽略，在步骤3导入数据后就正常了）
3. 登录`docker`启动的`mysql`，新建数据库`go-blog`,导入`go-blog/database/blog-mysql.sql`数据。
4. 访问url`http://127.0.0.0:8080`,后台url`http://127.0.0.0:8080/admin`,默认账户:`admin`,密码:`123456`

#### 方式三 源码安装

1. 把 Go Bbs 项目拉到本地 `git clone https://github.com/1920853199/go-blog.git`
2. 新建数据库，导入数据库文件，数据库文件/database/blog.sql
3. 修改项目配置信息

    ```
    # conf/app.conf

    appname = go-blog
    httpport = 8088
    runmode = dev
    EnableAdmin = false
    sessionon = true
    url = 127.0.0.1:8088

    [db]
    dbType = mysql
    dbUser = root
    dbPass = root
    dbHost = 127.0.0.1
    dbPort = 3306
    dbName = blog 
    ``` 
> 其他具体配置可以查看 beego 官网

4. 在 bo-bbs 根目录下执行`go run .` ，访问 http://127.0.0.1:8088 即可

> 后台访问链接 http://127.0.0.1:8088/admin, 默认账户:`admin`,密码:`123456`
