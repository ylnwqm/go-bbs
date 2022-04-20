package hotnews

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-bbs/models/admin"
	"go-bbs/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ImageData struct {
	B_64 string `json:"b_64"`
	Src  string `json:"url"`
}

func GetScreenshot(url string) {

	return
	param := strings.Split(url, ",")
	if len(param) < 2 {
		return
	}
	//json序列化
	postJson := `{"url":"` + param[1] + `"}`

	var jsonStr = []byte(postJson)

	req, err := http.NewRequest("POST", "http://152.32.212.226:16666/api", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("GetScreenshot Error : %s\n", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("GetScreenshot Error : %s\n", err.Error())
	}
	var img ImageData
	json.Unmarshal(body, &img)
	srcPath, err := utils.DownImage(img.Src)

	if err != nil {
		fmt.Printf("GetScreenshot Error : %s\n", err.Error())
	}

	hotId, err := strconv.Atoi(param[0])

	if err != nil {
		fmt.Printf("GetScreenshot Error : %s\n", err.Error())
	}
	fmt.Printf("Url:%s,Id:%d,ImgUrl:%s\n", img.Src, hotId, srcPath)
	err = admin.UpdateHotnewsById(&admin.Hotnews{
		Id:     int(hotId),
		ImgUrl: srcPath,
	})

	if err != nil {
		fmt.Printf("GetScreenshot Error : %s\n", err.Error())
	}

	return
}
