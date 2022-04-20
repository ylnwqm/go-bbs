package admin

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

var Msg = make(chan Notice, 10000)

const (
	SysMessage int = iota
	MsgMessage
	ChatMessage
)

var MsgTitle = map[int]string{SysMessage: "系统消息", MsgMessage: "新消息", ChatMessage: "您有一天新群聊消息"}

type Notice struct {
	Id        int    `orm:"column(id);auto" description:"ID"`
	SendId    int    `orm:"column(send_id)" description:"发送人"`
	ReceiveId int    `orm:"column(receive_id)" description:"接收人"`
	Type      int    `orm:"column(type)" description:"消息类型"`
	Title     string `orm:"column(title);size(50)" description:"标题"`
	Content   string `orm:"column(content);size(50)" description:"消息内容"`
	Status    int    `orm:"column(status);null" description:"1未读，2已读禁用"`
	Msg NoticeMessage `orm:"-"`
}

type NoticeMessage struct{
	Title string `json:"Title"`
	SendId int `json:"SendId"`
	ReceiveId int `json:"ReceiveId"`
	Username string `json:"Username"`
	Cover string `json:"Cover"`
	UserUrl string `json:"UserUrl"`
	Content string `json:"Content"`
	Replay string `json:"Replay"`
	Url string `json:"Url"`
	Date time.Time `json:"Date"`
}

func (t *Notice) TableName() string {
	return "notice"
}

func init() {
	orm.RegisterModel(new(Notice))
}

// AddNotice insert a new Notice into database and returns
// last inserted Id on success.
func AddNotice(m *Notice) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetNoticeById retrieves Notice by Id. Returns error if
// Id doesn't exist
func GetNoticeById(id int) (v *Notice, err error) {
	o := orm.NewOrm()
	v = &Notice{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetNoReadNiticeByReceiveId(receiveId int)(n []Notice,err error){
	o := orm.NewOrm()
	qs := o.QueryTable(new(Notice))
	var l []Notice
	_, err = qs.Filter("receive_id", receiveId).Filter("status", 1).OrderBy("-id").All(&l)

	if len(l) > 0 {
		for k,v := range l {
			l[k].Msg = UnmarshalNitice(v.Content)
		}
	}
	return l,err
}
func UnmarshalNitice(j string) NoticeMessage {
	var msg NoticeMessage
    json.Unmarshal([]byte(j), &msg)
	return msg
}
func GetNiticeByReceiveId(receiveId,limit,offset int)(n []Notice,err error){
	o := orm.NewOrm()
	qs := o.QueryTable(new(Notice))
	var l []Notice
	_, err = qs.Filter("receive_id", receiveId).OrderBy("-id").Limit(limit).Offset(offset).All(&l)
	if len(l) > 0 {
		for k,v := range l {
			l[k].Msg = UnmarshalNitice(v.Content)
		}
	}
	return l,err
}

// GetAllNotice retrieves all Notice matches certain condition. Returns empty list if
// no records exist
func GetAllNotice(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Notice))
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

	var l []Notice
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

// UpdateNotice updates Notice by Id and returns error if
// the record to be updated doesn't exist
func UpdateNoticeById(m *Notice) (err error) {
	o := orm.NewOrm()
	v := Notice{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteNotice deletes Notice by Id and returns error if
// the record to be deleted doesn't exist
func DeleteNotice(id int) (err error) {
	o := orm.NewOrm()
	v := Notice{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Notice{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func SetReadStatus(m *Notice,uid int) (err error) {
	o := orm.NewOrm()
	v := Notice{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		if v.ReceiveId != uid{
			return
		}
		v.Status = 2
		var num int64
		if num, err = o.Update(&v); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func GetNoticeCount(id int) (count int64, err error) {
	o := orm.NewOrm()

	count, err = o.QueryTable(new(Notice)).Filter("receive_id", id).Filter("status", 1).Count()

	return
}
