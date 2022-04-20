package api

import (
	"go-bbs/utils"
)

type AgentController struct {
	BaseController
}

func (ctl *AgentController) AgentUrl() {
	url := ctl.Ctx.Input.Param(":splat")
	url = utils.AesDecrypt(url, utils.Key)

	ctl.Redirect(url, 302)
}
