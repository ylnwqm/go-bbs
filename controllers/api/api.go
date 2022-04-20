package api

import (
    "bytes"
    "io/ioutil"
    "net/http"
	"encoding/json"
)

type ApiController struct {
	BaseController
}
type ImageData struct{
	B_64 string `json:"b_64"`
	Src  string `json:"url"`
}
func (c *ApiController) GetScreenshot(){
	response := make(map[string]interface{})

	url := c.GetString("url")

    //json序列化
    postJson := `{"url":"`+url+`"}`

    var jsonStr = []byte(postJson)
    

    req, err := http.NewRequest("POST", "http://152.32.212.226:16666/api", bytes.NewBuffer(jsonStr))
  
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
		response["code"] = 500
        response["msg"] = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
    }
    defer resp.Body.Close()

    
    body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response["code"] = 500
        response["msg"] = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
    var img ImageData
	json.Unmarshal(body,&img)

	response["Data"] = img
	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}
