package controllers

import (
	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
)

type MigrationController struct {
	beego.Controller
}

func (c *MigrationController) Get() {

	db := utils.ConnectDB()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Process{}, &models.ProcessData{})

	responseData := map[string]interface{}{
		"message": "Migration completed successfully!",
	}

	c.Data["json"] = responseData
	c.ServeJSON()
}
