package admin

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Hotnews struct {
	Id         int       `orm:"column(id);pk" description:"Id"`
	Category  *Category `orm:"rel(one)"`
	Title      string    `orm:"column(title);size(500)" description:"标题"`
	Url        string    `orm:"column(url);size(500)" description:"Url"`
	HotDesc    string    `orm:"column(hotDesc);size(500)" description:"热点描述"`
	ImgUrl     string    `orm:"column(imgUrl);size(250)" description:"图片连接"`
	Created    time.Time `orm:"column(created);type(datetime);null" description:"创建时间"`
	Rurl       string `orm:"-"`
	FMCreated  string `orm:"-"`
}

func (t *Hotnews) TableName() string {
	return "hotnews"
}

func init() {
	orm.RegisterModel(new(Hotnews))
}

// AddHotnews insert a new Hotnews into database and returns
// last inserted Id on success.
func AddHotnews(m *Hotnews) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetHotnewsById retrieves Hotnews by Id. Returns error if
// Id doesn't exist
func GetHotnewsById(id int) (v *Hotnews, err error) {
	o := orm.NewOrm()
	v = &Hotnews{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHotnews retrieves all Hotnews matches certain condition. Returns empty list if
// no records exist
func GetAllHotnews(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Hotnews))
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

	var l []Hotnews
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

// UpdateHotnews updates Hotnews by Id and returns error if
// the record to be updated doesn't exist
func UpdateHotnewsById(m *Hotnews) (err error) {
	o := orm.NewOrm()
	v := Hotnews{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		// fmt.Println(v)
		v.ImgUrl = m.ImgUrl
		if num, err = o.Update(&v); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func UpdateHotnewsByUrl(m *Hotnews) (err error) {
	o := orm.NewOrm()
	v := Hotnews{Url: m.Url}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		// fmt.Println(v)
		v.ImgUrl = m.ImgUrl
		if num, err = o.Update(&v); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}
// DeleteHotnews deletes Hotnews by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHotnews(id int) (err error) {
	o := orm.NewOrm()
	v := Hotnews{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Hotnews{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
