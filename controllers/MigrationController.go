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

	tokenString := c.Ctx.Input.Header("Authorization")
	_, role, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}

	if role != 1 {
		utils.CreateErrorResponse(&c.Controller, 400, "You are not authorized to perform this action.")
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
