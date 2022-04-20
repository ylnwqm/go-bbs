package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Books struct {
	Id            int       `orm:"column(id);auto" description:"主键"`
	EbookUrl      string    `orm:"column(ebook_url);size(500)"`
	Image         string    `orm:"column(image);size(500)"`
	Author        string    `orm:"column(author);size(500)"`
	Pubdate       string    `orm:"column(pubdate);size(500)"`
	Catalog       string    `orm:"column(catalog)"`
	Pages         string    `orm:"column(pages);size(500)"`
	Title         string    `orm:"column(title);size(500)"`
	AuthorIntro   string    `orm:"column(author_intro)" description:"作者简介"`
	Summary       string    `orm:"column(summary)" description:"简介"`
	MaxRaters     string    `orm:"column(max_raters);size(500)"`
	NumRaters     string    `orm:"column(num_raters);size(500)"`
	AverageRaters string    `orm:"column(average_raters);size(500)"`
	MinRaters     string    `orm:"column(min_raters);size(500)"`
	Isbn10        string    `orm:"column(Isbn10);size(500)"`
	Isbn13        string    `orm:"column(Isbn13);size(500)"`
	Ext           string    `orm:"column(ext)" description:"所有的信息"`
	CustomerId    int       `orm:"column(customer_id)" description:"用户ID"`
	Status        int       `orm:"column(status);null" description:"1在读，2想读，3以读，0删除"`
	Created       time.Time `orm:"column(created);type(datetime);null" description:"创建时间"`
}

func (t *Books) TableName() string {
	return "books"
}

func init() {
	orm.RegisterModel(new(Books))
}

// AddBooks insert a new Books into database and returns
// last inserted Id on success.
func AddBooks(m *Books) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetBooksById retrieves Books by Id. Returns error if
// Id doesn't exist
func GetBooksById(id int) (v *Books, err error) {
	o := orm.NewOrm()
	v = &Books{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBooks retrieves all Books matches certain condition. Returns empty list if
// no records exist
func GetAllBooks(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Books))
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

	var l []Books
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

// UpdateBooks updates Books by Id and returns error if
// the record to be updated doesn't exist
func UpdateBooksById(m *Books) (err error) {
	o := orm.NewOrm()
	v := Books{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBooks deletes Books by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBooks(id int) (err error) {
	o := orm.NewOrm()
	v := Books{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Books{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
