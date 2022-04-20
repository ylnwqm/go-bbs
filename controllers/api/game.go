package api



type GameController struct {
	BaseController
}

func (c *GameController) PacMan(){
	c.TplName = "home/game/index.html"
}
