package utils

import (
	"github.com/astaxie/beego"
)

func CreateErrorResponse(c *beego.Controller, statusCode int, message string) {
	c.Ctx.Output.SetStatus(statusCode)
	errorMessage := map[string]interface{}{
		"message": message,
	}

	c.Data["json"] = errorMessage
	c.ServeJSON()
}
