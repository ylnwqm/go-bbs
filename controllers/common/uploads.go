package common

import (
	"fmt"
	"go-bbs/utils"
	"math/rand"
	"path"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type UploadsController struct {
	beego.Controller
}

var extend = map[string]string{
	".jpg":  "jpg",
	".jpeg": "jpeg",
	".png":  "png",
	".gif":  "gif",
	".mp4":  "mp4",
	".avi":  "avi",
	".mov":  "mov",
}

func (c *UploadsController) Uploads() {

	response := make(map[string]interface{})

	f, h, err := c.GetFile("editormd-image-file")
	defer f.Close()
	if err != nil {
		response["message"] = err.Error()
		response["success"] = 0
	} else {

		ext := path.Ext(h.Filename)
		filename := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000)) + ext

		filepath := "static/uploads/images/" + time.Now().Format("20060102")
		if err = utils.CheckDir(filepath); err != nil {
			response["message"] = err.Error()
			response["success"] = 0
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
		path := filepath + "/" + filename
		err := c.SaveToFile("editormd-image-file", path)
		if err != nil {
			response["message"] = err.Error()
			response["success"] = 0
		} else {
			response["success"] = 1
			response["message"] = "Success."
			response["url"] = "/" + path
		}
	}

	c.Data["json"] = response
	c.ServeJSON()

}

func (c *UploadsController) UploadFiles() {

	response := make(map[string]interface{})

	f, h, err := c.GetFile("editormd-image-file")
	defer f.Close()
	if err != nil {
		response["message"] = err.Error()
		response["success"] = 0
	} else {

		//ext := path.Ext(h.Filename)
		//filename := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000)) + ext
		filepath := "static/uploads/files/" + time.Now().Format("20060102")
		if err = utils.CheckDir(filepath); err != nil {
			response["message"] = err.Error()
			response["success"] = 0
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
		path := filepath + "/" + h.Filename
		err := c.SaveToFile("editormd-image-file", path)
		if err != nil {
			response["message"] = err.Error()
			response["success"] = 0
		} else {
			response["success"] = 1
			response["message"] = "Success."
			response["url"] = "/" + path
			response["filename"] = h.Filename
		}
	}

	c.Data["json"] = response
	c.ServeJSON()

}

func (c *UploadsController) UploadsBbs() {

	response := make(map[string]interface{})

	f, h, err := c.GetFile("editormd-image-file")
	defer f.Close()
	if err != nil {
		response["message"] = err.Error()
		response["success"] = 0
	} else {

		ext := path.Ext(h.Filename)

		fmt.Println(ext)
		if _, ok := extend[ext]; !ok {
			response["message"] = "文件格式错误！"
			response["success"] = 0
		} else {
			filename := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000)) + ext

			filepath := "static/uploads/bbsfile/" + time.Now().Format("20060102")
			if err = utils.CheckDir(filepath); err != nil {
				response["message"] = err.Error()
				response["success"] = 0
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
			path := filepath + "/" + filename
			err := c.SaveToFile("editormd-image-file", path)
			if err != nil {
				response["message"] = err.Error()
				response["success"] = 0
			} else {
				response["ext"] = ext
				response["success"] = 1
				response["message"] = "Success."
				response["url"] = "/" + path
			}
		}
	}

	c.Data["json"] = response
	c.ServeJSON()

}
