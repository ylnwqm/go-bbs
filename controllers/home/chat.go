package home


type ChatController struct {
	BaseController
}

func (c *ChatController) Chat() {

	c.TplName = "home/" + c.Template + "/chat.html"
}


func (c *ChatController) ChatRoom() {

	c.TplName = "home/" + c.Template + "/chatroom.html"
}
