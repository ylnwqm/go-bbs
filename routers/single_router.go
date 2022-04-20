package routers

import (
	"github.com/astaxie/beego"

	"go-bbs/controllers/home"
)

func init() {

	beego.Router("/single_test", &home.SingleController{}, "Get:Test")
	beego.Router("/test_2", &home.SingleController{}, "Get:Test2")
	beego.Router("/tokyo-2020.html", &home.SingleController{}, "Get:Olympics")
	beego.Router("/google17a908a5fab67dec.html", &home.SingleController{}, "Get:Google17a908a5fab67dec")
}
