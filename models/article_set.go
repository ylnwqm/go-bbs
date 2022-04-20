package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ArticleSet struct {
	ArticleId  int       `orm:"column(article_id);null" description:"文章ID"`
	Created    time.Time `orm:"column(created);type(datetime);null" description:"创建时间"`
	CustomerId int       `orm:"column(customer_id);null" description:"用户ID"`
	Id         int       `orm:"column(id);auto" description:"主键"`
	Status     int       `orm:"column(status);null" description:"1可用，2禁用，3删除"`
	Title      string    `orm:"column(title);size(255);null" description:"名称"`
	Count      int64     `orm:"-"`
}

func (t *ArticleSet) TableName() string {
	return "article_set"
}

func init() {
	orm.RegisterModel(new(ArticleSet))
}

// AddArticleSet insert a new ArticleSet into database and returns
// last inserted Id on success.
func AddArticleSet(m *ArticleSet) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetArticleSetById retrieves ArticleSet by Id. Returns error if
// Id doesn't exist
func GetArticleSetById(id int) (v *ArticleSet, err error) {
	o := orm.NewOrm()
	v = &ArticleSet{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetArticleSetCustomerId(CustomerId int) (v []*ArticleSet, err error) {
	var l []*ArticleSet
	o := orm.NewOrm()
	qs := o.QueryTable(new(ArticleSet))
	qs = qs.Filter("customer_id", CustomerId)
	qs = qs.Filter("status", 1)
	_, err = qs.Limit(500, 0).GroupBy("title").All(&l)
	return l, err
}

func GetArticleSetByArticleId(ArticleId, CustomerId int) (v ArticleSet, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ArticleSet))
	qs = qs.Filter("customer_id", CustomerId)
	qs = qs.Filter("status", 1)
	qs = qs.Filter("article_id", ArticleId)
	// v = &ArticleSet{ArticleId: ArticleId, Status: 1, CustomerId: CustomerId}
	err = qs.One(&v)
	return v, err
}

// GetAllArticleSet retrieves all ArticleSet matches certain condition. Returns empty list if
// no records exist
func GetAllArticleSet(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ArticleSet))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []ArticleSet
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateArticleSet updates ArticleSet by Id and returns error if
// the record to be updated doesn't exist
func UpdateArticleSetById(m *ArticleSet) (err error) {
	o := orm.NewOrm()
	v := ArticleSet{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteArticleSet deletes ArticleSet by Id and returns error if
// the record to be deleted doesn't exist
func DeleteArticleSet(id int) (err error) {
	o := orm.NewOrm()
	v := ArticleSet{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ArticleSet{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
