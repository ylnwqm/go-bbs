package utils

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/require"
)

func TestUtilsGetTopic(t *testing.T) {
	topic, content := GetTopic("#asd#000")
	require.EqualValues(t, topic, "#asd#")
	require.EqualValues(t, content, "000")

	topic, content = GetTopic("#你好#000")
	require.EqualValues(t, topic, "#你好#")
	require.EqualValues(t, content, "000")

	topic, content = GetTopic("你#好#000")
	require.EqualValues(t, topic, "")
	require.EqualValues(t, content, "你#好#000")

	topic, content = GetTopic("#你好#000#")
	require.EqualValues(t, topic, "#你好#")
	require.EqualValues(t, content, "000#")

	content = GetTopicContent("#你好#000#")
	require.EqualValues(t, content, "#你好#")
	
	content = GetTopicContent("霓虹")
	require.EqualValues(t, content, "#霓虹#")

	content = GetTopicContent("#你好")
	fmt.Println(content)
	require.EqualValues(t, content, "#你好#")
}

func TestUtilsDownImage(t *testing.T) {

	fmt.Println(beego.WorkPath)
	path, err := DownImage("https://img9.doubanio.com/view/thing_review/l/public/p6153454.webp")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(path)
	require.EqualValues(t, path, "/static/uploads/images/202106/p6153454.webp")
}
