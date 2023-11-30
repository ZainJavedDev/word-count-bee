package controllers

import (
	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
)

type MigrationController struct {
	beego.Controller
}

func (c *MigrationController) Post() {

	dbKey := c.Ctx.Input.Header("Key")

	envDBKey := utils.GoDotEnvVariable("DB_ACCESS_KEY")

	if dbKey != envDBKey {
		utils.CreateErrorResponse(&c.Controller, 422, "Invalid key")
	}

	db := utils.ConnectDB()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Process{}, &models.ProcessData{})

	responseData := map[string]interface{}{
		"message": "Migration completed successfully!",
	}

	c.Data["json"] = responseData
	c.ServeJSON()
}
