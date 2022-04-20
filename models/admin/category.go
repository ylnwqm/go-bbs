package admin

import (
	"github.com/astaxie/beego/orm"
)

type Category struct {
	Id     int
	Name   string
	Pid    int
	Sort   int
	Status int
	Description string
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(Category))
}
