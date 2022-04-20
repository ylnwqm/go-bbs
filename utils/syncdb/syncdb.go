package syncdb

import (
	"database/sql"
	"fmt"
	"go-bbs/models/admin"
	model "go-bbs/models/admin"
	"go-bbs/utils"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var o orm.Ormer

func Syncdb() {
	createdb()
	Connect()
	o = orm.NewOrm()
	// 数据库别名
	name := "default"
	// drop table 后再建表
	force := true
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}

	insertUser()
	insertSetting()
	insertMenu()
	insertArticle()
	fmt.Println("database init is complete.\nPlease restart the application")

}

//数据库连接
func Connect() {
	var dsn string
	db_type := beego.AppConfig.String("db::dbType")
	db_host := beego.AppConfig.String("db::dbHost")
	db_port := beego.AppConfig.String("db::dbPort")
	db_user := beego.AppConfig.String("db::dbUser")
	db_pass := beego.AppConfig.String("db::dbPass")
	db_name := beego.AppConfig.String("db::dbName")
	switch db_type {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
		break
	// case "postgres":
	// 	orm.RegisterDriver("postgres", orm.DRPostgres)
	// 	dsn = fmt.Sprintf("dbname=%s host=%s  user=%s  password=%s  port=%s  sslmode=%s", db_name, db_host, db_user, db_pass, db_port, db_sslmode)
	// case "sqlite3":
	// 	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	// 	if db_path == "" {
	// 		db_path = "./"
	// 	}
	// 	dsn = fmt.Sprintf("%s%s.db", db_path, db_name)
	// 	break
	default:
		beego.Critical("Database driver is not allowed:", db_type)
	}
	orm.RegisterDataBase("default", db_type, dsn)
}

//创建数据库
func createdb() {

	db_type := beego.AppConfig.String("db::dbType")
	db_host := beego.AppConfig.String("db::dbHost")
	db_port := beego.AppConfig.String("db::dbPort")
	db_user := beego.AppConfig.String("db::dbUser")
	db_pass := beego.AppConfig.String("db::dbPass")
	db_name := beego.AppConfig.String("db::dbName")

	var dsn string
	var sqlstring string
	switch db_type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", db_user, db_pass, db_host, db_port)
		sqlstring = fmt.Sprintf("CREATE DATABASE  if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)
		break
	// case "postgres":
	// 	dsn = fmt.Sprintf("host=%s  user=%s  password=%s  port=%s  sslmode=%s", db_host, db_user, db_pass, db_port, db_sslmode)
	// 	sqlstring = fmt.Sprintf("CREATE DATABASE %s", db_name)
	// 	break
	// case "sqlite3":
	// 	if db_path == "" {
	// 		db_path = "./"
	// 	}
	// 	dsn = fmt.Sprintf("%s%s.db", db_path, db_name)
	// 	os.Remove(dsn)
	// 	sqlstring = "create table init (n varchar(32));drop table init;"
	// 	break
	default:
		beego.Critical("Database driver is not allowed:", db_type)
	}

	db, err := sql.Open(db_type, dsn)
	if err != nil {
		panic(err.Error())
	}
	r, err := db.Exec(sqlstring)
	if err != nil {
		log.Println(err)
		log.Println(r)
	} else {
		log.Println("Database ", db_name, " created")
	}
	defer db.Close()

}

func insertUser() {
	fmt.Println("insert admin user ...")
	u := new(model.User)
	u.Name = "admin"
	u.Password = utils.PasswordMD5("123456", "admin")
	u.Email = "go-bbs@gmail.com"
	u.Created = time.Now()
	u.Status = 1

	o = orm.NewOrm()
	o.Insert(u)
	o.Insert(&admin.Customer{
		Uid:      "1",
		Username: "User",
		Email:    "go-bbs@gmail.com",
		Password: utils.PasswordMD5("123456", "nihongdengxia.com"),
		Status:   1,
		Image:    "/static/images/cover_default.jpg",
		Integral: 20,
		Fans:     0,
		Focus:    0,
		Created:  time.Now(),
		Updated:  time.Now(),
	},
	)
	fmt.Println("insert user end")
}

func insertSetting() {
	fmt.Println("insert setting ...")

	o = orm.NewOrm()

	err := o.Begin()

	vw := beego.AppConfig.String("view")
	if vw == "" {
		vw = "leechan"
	}

	settings := []model.Setting{
		{Name: "title", Value: "GO-BBS - Golang开源社区系统"},
		{Name: "template", Value: vw},
		{Name: "mobile", Value: vw},
		{Name: "limit", Value: "20"},
		{Name: "tag", Value: "Go社区系统,Go开源社区系统,go-bbs,Go,bbs,Go社区,bbsgo"},
		{Name: "image", Value: "/static/img/logo.jpg"},
		{Name: "keyword", Value: "Go社区系统,Go开源社区系统,go-bbs,Go,bbs,Go社区,bbsgo"},
		{Name: "description", Value: "Golang最优秀最简洁最好看的社区开源系统之一"},
		{Name: "notice", Value: "欢迎使用Golang开源社区系统GO-BBS！"},
		{Name: "notice_html_code", Value: "欢迎使用Golang开源社区系统GO-BBS！"},
		{Name: "notice_markdown_doc", Value: "欢迎使用Golang开源社区系统GO-BBS！"},
	}

	_, err = o.InsertMulti(11, settings)

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	fmt.Println("insert setting end")
}

func insertArticle() {
	fmt.Println("insert Article ...")

	o = orm.NewOrm()

	err := o.Begin()
	id, err := o.Insert(&admin.Category{
		Name:   "默认分类",
		Sort:   100,
		Pid:    0,
		Status: 1,
	},
	)

	article := []model.Article{
		admin.Article{
			Title:    "关于 Go Bbs 社区系统",
			Tag:      "Go社区系统,Go开源社区系统,go-bbs,Go,bbs,Go社区,bbsgo",
			Desc:     "Go Bbs 是一个基于 Beego 开发的可切换模板的 BBS 社交博客系统。<br>它功能齐全，安装简单便捷，一个`go` 命令就可以跑起来了，并且没有什么第三方依赖。页面简介优美。前端是HTML+JS+CSS，不需要掌握一些前端技术栈也能轻松自定义页面。目前是 go bbs 出生的第三年了，他前身是 go-blog 。<br><br> 现在我们正式把 go bbs 官网以及文档搭建起来，文档后面会持续补充完整。在这里也期待大家的支持。<p><a href='http://go-bbs.com/' title='GO BBS 官网'>GO BBS 官网</a></p>",
			Html:     "<p>Go Bbs 是一个基于 Beego 开发的可切换模板的 BBS 社交博客系统。<br>它功能齐全，安装简单便捷，一个<code>go</code> 命令就可以跑起来了，并且没有什么第三方依赖。页面简介优美。前端是HTML+JS+CSS，不需要掌握一些前端技术栈也能轻松自定义页面。目前是 go bbs 出生的第三年了，他前身是 go-blog 。</p><p>现在我们正式把 go bbs 官网以及文档搭建起来，文档后面会持续补充完整。在这里也期待大家的支持。</p><p><a href='http://go-bbs.com/' title='GO BBS 官网'>GO BBS 官网</a></p>",
			Remark:   "Go Bbs 是一个基于 Beego 开发的可切换模板的 BBS 社交博客系统。它功能齐全，安装简单便捷，一个`go` 命令就可以跑起来了，并且没有什么第三方依赖。页面简介优美。前端是HTML+JS+CSS，不需要掌握一些前端技术栈也能轻松自定义页面",
			Url:      "",
			Cover:    "",
			Status:   1,
			Customer: &admin.Customer{Id: 1},
			Category: &admin.Category{int(id), "", 0, 0, 0, ""},
		},
	}

	_, err = o.InsertMulti(11, article)

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	fmt.Println("insert Article end")
}

func insertMenu() {
	fmt.Println("insert menu ...")

	o = orm.NewOrm()

	err := o.Begin()

	menus := []model.Menu{
		{Title: "首页", Url: "/", Sort: 100, Pid: 0},
		//{Title: "社区", Url: "/bbs.html", Sort: 100, Pid: 0},
		//{Title: "话题广场", Url: "/square.html", Sort: 100, Pid: 0},
		{Title: "文章", Url: "/articles.html", Sort: 100, Pid: 0},
		//{Title: "友情链接", Url: "/links.html", Sort: 100, Pid: 0},
		//{Title: "会员", Url: "/members", Sort: 100, Pid: 0},
	}

	_, err = o.InsertMulti(11, menus)

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	fmt.Println("insert menu end")
}
