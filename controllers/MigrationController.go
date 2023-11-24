package controllers

import (
	"encoding/json"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
)

type MigrationController struct {
	beego.Controller
}

func (c *MigrationController) Get() {

	db := utils.ConnectDB()
	db.AutoMigrate(&models.Process{}, &models.ProcessData{})

	responseData := map[string]interface{}{
		"message": "Migration completed successfully!",
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	c.Ctx.Output.Body(jsonData)
}
