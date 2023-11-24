package utils

import (
	"encoding/json"
	"log"

	"github.com/astaxie/beego"
)

func CreateErrorResponse(c *beego.Controller, statusCode int, message string) {
	c.Ctx.Output.SetStatus(statusCode)
	errorMessage := map[string]interface{}{
		"message": message,
	}
	jsonData, err := json.Marshal(errorMessage)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		log.Fatal(err)
	}
	c.Ctx.Output.Body(jsonData)
}
