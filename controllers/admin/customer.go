package admin

import (
	"encoding/json"
	models "go-bbs/models/admin"
	"go-bbs/utils"
	"strconv"

	"github.com/astaxie/beego/orm"
)

// CustomerController operations for Customer
type CustomerController struct {
	BaseController
}

// URLMapping ...
func (c *CustomerController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Customer
// @Param	body		body 	models.Customer	true		"body for Customer content"
// @Success 201 {int} models.Customer
// @Failure 403 body is empty
// @router /customer [post]
func (c *CustomerController) Post() {
	var v models.Customer
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddCustomer(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Customer by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Customer
// @Failure 403 :id is empty
// @router /customer/:id [get]
func (c *CustomerController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetCustomerById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Customer
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Customer
// @Failure 403
// @router /customer/ [get]
func (c *CustomerController) GetAll() {

	limit := int64(15)               // 一页的数量
	page, _ := c.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit     // 偏移量

	name := c.GetString("name")
	status, _ := c.GetInt("status", 0)

	c.Data["Status"] = status
	c.Data["Name"] = name

	o := orm.NewOrm()

	var users []*models.Customer
	qs := o.QueryTable(new(models.Customer))

	// 状态
	if status != 0 {
		qs = qs.Filter("status", status)
	}

	// 名称
	if name != "" {
		qs = qs.Filter("name__icontains", name)
	}

	qs = qs.Filter("email__isnull", false)

	// 获取数据
	_, err := qs.OrderBy("-id").Limit(limit).Offset(offset).All(&users)

	if err != nil {
		c.Abort("404")
	}

	// 统计
	count, err := qs.Count()
	if err != nil {
		c.Abort("404")
	}

	c.Data["Data"] = &users
	c.Data["Paginator"] = utils.GenPaginator(page, limit, count)
	c.Data["StatusText"] = models.Status
	c.TplName = "admin/customer-list.html"

}

func (c *CustomerController) Status() {
	id, _ := c.GetInt("id", 0)
	status, _ := c.GetInt("status", 1)

	response := make(map[string]interface{})

	o := orm.NewOrm()
	customer := models.Customer{Id: id}
	if o.Read(&customer) == nil {
		if status == 1 {
			status = 2
		} else {
			status = 1
		}
		customer.Status = status

		if _, err := o.Update(&customer); err == nil {
			response["msg"] = "禁用成功！"
			response["code"] = 200
			response["id"] = id
		} else {
			response["msg"] = "禁用失败！"
			response["code"] = 500
			response["err"] = err.Error()
		}
	} else {
		response["msg"] = "禁用失败！"
		response["code"] = 500
		response["err"] = "ID 不能为空！"
	}

	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

// Put ...
// @Title Put
// @Description update the Customer
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Customer	true		"body for Customer content"
// @Success 200 {object} models.Customer
// @Failure 403 :id is not int
// @router /customer/:id [put]
func (c *CustomerController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Customer{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateCustomerById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Customer
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /customer/:id [delete]
func (c *CustomerController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteCustomer(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
