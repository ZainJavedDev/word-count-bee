package controllers

import (
	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
)

type RoleChangeController struct {
	beego.Controller
}

type RoleChangeUser struct {
	Username string `form:"username"`
}

func (c *RoleChangeController) Post() {
	dbKey := c.Ctx.Input.Header("Key")

	envDBKey := utils.GoDotEnvVariable("DB_ACCESS_KEY")

	if dbKey != envDBKey {
		utils.CreateErrorResponse(&c.Controller, 422, "Invalid key")
	}

	db := utils.ConnectDB()
	defer db.Close()

	var roleChangeUser RoleChangeUser
	err := c.ParseForm(&roleChangeUser)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid form data.")
	}

	var user models.User

	err = db.Where("username = ?", roleChangeUser.Username).Find(&user).Error
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 404, err.Error())
		return
	}

	if user.Role == 0 {
		user.Role = 1
	} else {
		user.Role = 0
	}

	db.Save(&user)

	responseData := map[string]interface{}{
		"message": "Role changed successfully!",
	}

	c.Data["json"] = responseData
	c.ServeJSON()

}
