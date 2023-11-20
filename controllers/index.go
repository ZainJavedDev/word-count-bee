package controllers

import (
	"encoding/json"

	"github.com/MrNi8mare/word-count-bee/utils"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type Message struct {
	FilePath string `json:"filepath"`
	Routines int    `json:"routines"`
}

func (c *MainController) Post() {

	var message Message
	decoder := json.NewDecoder(c.Ctx.Request.Body)
	err := decoder.Decode(&message)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	totalCounts, routines, timeTaken := utils.ProcessFile(message.FilePath, message.Routines)
	timeTakenString := timeTaken.String()

	responseData := map[string]interface{}{
		"totalCounts": totalCounts,
		"routines":    routines,
		"timeTaken":   timeTakenString,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	c.Ctx.Output.Body(jsonData)
}
