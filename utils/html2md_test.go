package utils

import (
	"fmt"
	"testing"
)

func TestHtmlRemoteImg2LocalImg(t *testing.T) {

	html := `<body>

				<div id="div1">DIV1</div>
				<div>DIV2</div>
				<span>SPAN</span>
				<img src="https://img2.doubanio.com/view/thing_review/l/public/p2134132.webp" width="670">
			</body>`

	str := HtmlRemoteImg2LocalImg(html)
	fmt.Println(str)
}
