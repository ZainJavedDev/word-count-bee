package controllers

import (
	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
)

type ProcessController struct {
	beego.Controller
}

func (c *ProcessController) Post() {

	tokenString := c.Ctx.Input.Header("Authorization")
	userID, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}
	db := utils.ConnectDB()
	defer db.Close()
	userProcesses := []models.Process{}
	db.Where("user_id = ?", userID).Preload("User").Find(&userProcesses)

	c.Data["json"] = userProcesses
	c.ServeJSON()

}
