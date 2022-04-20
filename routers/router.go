package routers

import (
	"go-bbs/controllers/admin"
	"go-bbs/controllers/common"
	"go-bbs/controllers/home"
	"go-bbs/controllers/installer"
	"go-bbs/controllers/websocket"
	"go-bbs/controllers/wechat"

	"go-bbs/controllers/api"

	"github.com/astaxie/beego"
)

func init() {

	// 安装功能
	// beego.InsertFilter("*",beego.BeforeExec,filter.Installer)

	adminNs := beego.NewNamespace("/admin",
		beego.NSInclude(
			&admin.MenuController{},
			&admin.LinkController{},
			&admin.FileController{},
			&admin.CronController{},
			&admin.CustomerController{},
			&admin.AdController{},
		),

		beego.NSRouter("/ad-show", &admin.AdController{}, "get:Add"),
		beego.NSRouter("/ad-edit", &admin.AdController{}, "get:Edit"),
		beego.NSRouter("/group/delete", &admin.AdController{}, "delete:DeleteGroup"),
		beego.NSRouter("/get-review", &admin.CronController{}, "get:GetReview"),

		beego.NSRouter("/user", &admin.UserController{}, "get:List;post:Save"),
		beego.NSRouter("/user/edit", &admin.UserController{}, "get:Put"),
		beego.NSRouter("/user/status", &admin.UserController{}, "Post:Delete"),
		beego.NSRouter("/user/update", &admin.UserController{}, "Post:Update"),
		beego.NSRouter("/user/add", &admin.UserController{}, "get:Add"),

		beego.NSRouter("/customer/status", &admin.CustomerController{}, "Post:Status"),

		// 站点设置
		beego.NSRouter("/setting", &admin.SettingController{}, "get:Add"),
		beego.NSRouter("/setting/save", &admin.SettingController{}, "post:Save"),
		beego.NSRouter("/notice", &admin.SettingController{}, "get:Notice"),
		beego.NSRouter("/notice/save", &admin.SettingController{}, "post:NoticeSave"),
		beego.NSRouter("/about", &admin.SettingController{}, "get:About"),
		beego.NSRouter("/about/save", &admin.SettingController{}, "post:AboutSave"),
		beego.NSRouter("/contact", &admin.SettingController{}, "get:Contact"),
		beego.NSRouter("/contact/save", &admin.SettingController{}, "post:ContactSave"),

		// 后台文章模块
		beego.NSRouter("/welcome", &admin.MainController{}, "get:Welcome"),
		beego.NSRouter("/article", &admin.ArticleController{}, "get:List;post:Save"),
		beego.NSRouter("/article/edit", &admin.ArticleController{}, "get:Put"),
		beego.NSRouter("/article/delete", &admin.ArticleController{}, "Post:Delete"),
		beego.NSRouter("/article/update", &admin.ArticleController{}, "Post:Update"),
		beego.NSRouter("/article/add", &admin.ArticleController{}, "get:Add"),
		beego.NSRouter("/article/top", &admin.ArticleController{}, "Post:Top"),
		beego.NSRouter("/article/get", &admin.ArticleResourcesController{}, "Post:GetArticle"),

		beego.NSRouter("/article/cron/get", &admin.ArticleResourcesController{}, "Get:GetCron"),

		// 后台分类模块
		beego.NSRouter("/cate", &admin.CateController{}, "get:List;post:Save"),
		beego.NSRouter("/cate/add", &admin.CateController{}, "get:Add"),
		beego.NSRouter("/cate/edit", &admin.CateController{}, "get:Put"),
		beego.NSRouter("/cate/delete", &admin.CateController{}, "Post:Delete"),
		beego.NSRouter("/cate/update", &admin.CateController{}, "Post:Update"),
		// 后台登录
		beego.NSRouter("/login", &admin.LoginController{}, "Get:Sign;Post:Login"),
		beego.NSRouter("/logout", &admin.LoginController{}, "Get:Logout"),

		// 后台评论
		beego.NSRouter("/review", &admin.ReviewController{}, "get:List"),
		beego.NSRouter("/review/edit", &admin.ReviewController{}, "Get:Put"),
		beego.NSRouter("/review/delete", &admin.ReviewController{}, "Post:Delete"),
		beego.NSRouter("/review/update", &admin.ReviewController{}, "Post:Update"),

		// 后台留言
		beego.NSRouter("/message", &admin.MessageController{}, "get:List"),

		// 后台留言模块
		beego.NSRouter("/message/update", &admin.MessageController{}, "Post:Update"),
		beego.NSRouter("/message/edit", &admin.MessageController{}, "Get:Put"),
		beego.NSRouter("/message/delete", &admin.MessageController{}, "Post:Delete"),

		beego.NSRouter("/topic", &admin.TopicController{}, "get:List"),
		beego.NSRouter("/topic/edit", &admin.TopicController{}, "get:Put"),
		beego.NSRouter("/topic/update", &admin.TopicController{}, "post:Update"),

		beego.NSRouter("/logs", &admin.LogController{}, "get:GetAll"),
	)

	beego.AddNamespace(adminNs)
	// 公众号
	beego.Router("/wechat", &wechat.MainController{}, "Get:CheckToken;Post:Hello")
	beego.Router("/wechat/create/menu", &wechat.MenuController{}, "Get:CreateMenu;Post:CreateMenu")
	beego.Router("/wechat/user/get", &wechat.UserController{}, "Get:GetUser")
	beego.Router("/wechat/user/list", &wechat.UserController{}, "Get:List")
	beego.Router("/wechat/addnews", &wechat.MaterialController{}, "Get:AddNews")

	beego.Router("/", &home.MainController{})
	beego.Router("/admin", &admin.MainController{}, "get:Index")
	// 前台列表
	beego.Router("/list.html", &home.ArticleController{}, "get:List")
	beego.Router("/article/category/:id([0-9]+).html", &home.ArticleController{}, "get:List")
	beego.Router("/news.html", &home.ArticleController{}, "Get:News")
	// 前台详情
	beego.Router("/detail/:id([0-9]+).html", &home.ArticleController{}, "get:Detail")
	// 前台统计文章PV
	beego.Router("/pv/:id([0-9]+).html", &home.ArticleController{}, "get:Pv")
	// 前台留言列表
	beego.Router("/message.html", &home.MessageController{}, "get:Get")
	// 前台保存
	beego.Router("/message/save", &home.MessageController{}, "Post:Save")

	// 评论保存
	beego.Router("article/review", &home.ArticleController{}, "Post:Review")
	beego.Router("article/review/:id([0-9]+).html", &home.ArticleController{}, "Get:ReviewList")
	beego.Router("article/like", &home.ArticleController{}, "Post:Like")

	// 文件上传
	beego.Router("/uploads.html", &common.UploadsController{}, "Post:Uploads")         // 图片
	beego.Router("/uploadfiles.html", &common.UploadsController{}, "Post:UploadFiles") // 文件
	beego.Router("/uploadbbs.html", &common.UploadsController{}, "Post:UploadsBbs")    // 文件

	// 安装
	beego.Router("/installer", &installer.InstallController{}, "Get:CheckEnv")
	beego.Router("/installer/create", &installer.InstallController{}, "Get:Install")

	beego.Router("/login.html", &home.LoginController{}, "Get:Sign")
	beego.Router("/login.html", &home.LoginController{}, "Post:Login")
	beego.Router("/logout", &home.LoginController{}, "Get:Logout")
	beego.Router("/get-captcha.html", &home.LoginController{}, "get:GetCaptcha")
	beego.Router("/reg.html", &home.LoginController{}, "Post:Regist")

	beego.Router("/bbs.html", &home.BbsController{}, "Get:Bbs")
	beego.Router("/bbs/view/switch/*", &home.BbsController{}, "Get:SaveView")
	beego.Router("/bbs/:id([0-9]+).html", &home.BbsController{}, "Get:Bbs")

	beego.Router("/topic/detail/:id([0-9]+).html", &home.BbsController{}, "Get:Topic")
	beego.Router("/bbs/save", &home.BbsController{}, "Post:Save")
	beego.Router("/bbs/review/save", &home.BbsController{}, "Post:SaveReview")
	beego.Router("/bbs/like", &home.BbsController{}, "Post:Like")
	beego.Router("/bbs/detail/:id([0-9]+).html", &home.BbsController{}, "Get:BbsDetail")
	beego.Router("/bbs/create", &home.BbsController{}, "Get:BbsCreate")
	beego.Router("/images.html", &home.ImagesController{}, "get:List")
	beego.Router("/square", &home.TopicController{}, "Get:Get")

	beego.Router("/profile.html", &home.CustomerController{}, "Get:Profile")
	beego.Router("/profile.html", &home.CustomerController{}, "Post:Put")
	beego.Router("/user/:id([0-9]+)", &home.CustomerController{}, "Get:UserInfo")
	beego.Router("/members", &home.CustomerController{}, "Get:List")

	beego.Router("/articles.html", &home.ArticleController{}, "get:List")
	//beego.Router("/books.html", &home.BooksController{}, "get:List")
	beego.Router("/books.html", &home.BooksController{}, "get:Books")
	beego.Router("/search/books", &home.BooksController{}, "get:SearchBooks")
	beego.Router("/select/books", &home.BooksController{}, "get:SelectBooks")
	beego.Router("/save/books", &home.BooksController{}, "post:SaveBooks")
	beego.Router("/setstatus/books", &home.BooksController{}, "post:SetStatus")
	beego.Router("/like/books", &home.BooksController{}, "post:SetLike")

	// 关注
	beego.Router("/user/fans.html", &home.CustomerController{}, "Post:Fans")

	beego.Router("/about.html", &home.SingleController{}, "get:About")
	beego.Router("/contact.html", &home.SingleController{}, "get:Contact")

	beego.Router("/quotes.html", &home.SingleController{}, "get:Quotes")

	beego.Router("/music.html", &common.MusicController{}, "get:Get")
	beego.Router("/music/update", &common.MusicController{}, "get:Update")
	beego.Router("/sitemap.xml", &home.SingleController{}, "get:Sitemap")
	beego.Router("/links.html", &home.SingleController{}, "get:Links")
	beego.Router("/time.html", &home.SingleController{}, "get:History")

	beego.Router("/notice.html", &home.NoticeController{}, "get:Notice")
	beego.Router("/notice/send", &home.NoticeController{}, "get:Send")

	beego.Router("/hotnews", &home.HotnewsController{}, "get:Get")
	beego.Router("/get/danmu", &home.HotnewsController{}, "get:GetInc")
	beego.Router("/hotnews/history", &home.HotnewsController{}, "get:GetHistoryHot")

	beego.Router("/hotnews/history/page", &home.HotnewsController{}, "get:GetHistoryHotByDate")
	beego.Router("/hotnews/detail/:id([0-9]+).html", &home.HotnewsController{}, "get:HotDetail")

	beego.Router("/hot", &home.HotnewsController{}, "get:List")

	beego.Router("/chat.html", &home.ChatController{}, "get:Chat")
	beego.Router("/chatroom.html", &home.ChatController{}, "get:ChatRoom")

	beego.Router("/api/v1/download/*", &home.SingleController{}, "get:Download")
	beego.Router("/search", &home.SearchController{}, "get:Search")
	beego.Router("/articles/create", &home.CustomerArticleController{}, "get:Create")
	beego.Router("/articles/create", &home.CustomerArticleController{}, "post:Save")
	beego.Router("/articles/edit", &home.CustomerArticleController{}, "get:Put")
	beego.Router("/articles/edit", &home.CustomerArticleController{}, "post:Edit")
	beego.Router("/articleset", &home.ArticleSetController{}, "Get:GetArticleBySet")

	beego.Router("/games/pacman.html", &api.GameController{}, "get:PacMan")

	beego.Router("/dailyhot/detail/:id([0-9]+).html", &home.SingleController{}, "get:DailyhotDetail")

	beego.Router("/danmu/switch/*", &home.SingleController{}, "Get:SwitchDanmu")
	api := beego.NewNamespace("/api/v1",

		beego.NSRouter("/bbs/detail/:id([0-9]+).html", &api.BbsController{}, "Get:Bbs"),
		beego.NSRouter("/createsitemap0x001.html", &api.ToolController{}, "Get:CSitemap"),
		beego.NSRouter("/to/*", &home.SingleController{}, "Get:AgentUrl"),
		beego.NSRouter("/getScreenshot", &api.ApiController{}, "Get:GetScreenshot"),
		beego.NSRouter("/daily", &home.BbsController{}, "Post:SaveDaily"),
	)

	beego.AddNamespace(api)

	beego.Router("/ws", &websocket.NoticeController{}, "get:Join")
}
