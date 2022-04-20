package utils

import (
	"bufio"
	"fmt"
	models "go-bbs/models/admin"
	"os"
	"strings"
)

const ROUTER_PATH = "./routers/single_router.go"

func CreateRouters() error {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 1000
	var offset int64 = 0

	query["status"] = "1"

	l, err := models.GetAllRouter(query, fields, sortby, order, offset, limit)

	if err != nil {
		return err
	}

	routers := `
package routers

import (
	"github.com/astaxie/beego"
`
	var ims = make(map[string]bool)
	var path string

	for _, r := range l {
		v := r.(models.Router)
		ip := getBetweenStr(v.Controller, "&", ".")
		if _, ok := ims[ip]; !ok {
			ims[ip] = true
		}

		path += fmt.Sprintf(`
	beego.Router("%s", %s, "%s")`, v.Path, v.Controller, v.Methods)
	}

	for im := range ims {
		routers += fmt.Sprintf(`
	"go-bbs/controllers/%s"`, im)
	}
	routers += `
)	
`
	routers += `
func init(){
`
	routers += path + `
}	
`
	outputFile, outputError := os.OpenFile(ROUTER_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation，%s\n", outputError.Error())
		return outputError
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	outputWriter.WriteString(routers)
	outputWriter.Flush()
	return nil
}

func getBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	// 需要把签名的&去掉因此[1:m]
	str = string([]byte(str)[1:m])
	return str
}
