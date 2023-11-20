package controllers

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/MrNi8mare/word-count-bee/utils"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type Message struct {
	FilePath string `form:"file"`
	Routines int    `form:"routines"`
}

func (c *MainController) Post() {

	var message Message
	err := c.ParseForm(&message)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	uploadedFile, header, err := c.GetFile("file")
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}
	defer uploadedFile.Close()

	uploadDir := "./uploads/"
	err = os.MkdirAll(uploadDir, 0777)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	filePath := filepath.Join(uploadDir, header.Filename)

	outputFile, err := os.Create(filePath)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, uploadedFile)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	totalCounts, routines, timeTaken := utils.ProcessFile(filePath, message.Routines)
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
