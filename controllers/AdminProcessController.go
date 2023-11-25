package controllers

import (
	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
)

type AdminProcessController struct {
	beego.Controller
}

func (c *AdminProcessController) Post() {

	tokenString := c.Ctx.Input.Header("Authorization")
	_, role, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}

	if role != 1 {
		utils.CreateErrorResponse(&c.Controller, 400, "You are not authorized to perform this action.")
	}

	db := utils.ConnectDB()
	defer db.Close()
	allProcesses := []models.Process{}
	db.Preload("ProcessData").Find(&allProcesses)

	c.Data["json"] = allProcesses
	c.ServeJSON()

}
