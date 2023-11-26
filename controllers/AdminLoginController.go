package controllers

import (
	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/astaxie/beego"
)

type AdminLoginController struct {
	beego.Controller
}

func (c *AdminLoginController) Post() {

	var loginData LoginData
	err := c.ParseForm(&loginData)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	if loginData.Username == "" || loginData.Password == "" {
		errorMessage := "Username and password are required"
		utils.CreateErrorResponse(&c.Controller, 400, errorMessage)
	}

	db := utils.ConnectDB()
	defer db.Close()
	var userFromDB models.User
	result := db.Where("username = ?", loginData.Username).First(&userFromDB)
	if result.Error != nil {
		errorMessage := "Invalid credentials. Please check your username and password."
		utils.CreateErrorResponse(&c.Controller, 401, errorMessage)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(loginData.Password))
	if err != nil {
		errorMessage := "Invalid credentials. Please check your username and password."
		utils.CreateErrorResponse(&c.Controller, 401, errorMessage)
	}

	if userFromDB.Role != 1 {
		utils.CreateErrorResponse(&c.Controller, 401, "You are not authorized to perform this action.")
	}

	c.Ctx.Output.SetStatus(200)
	tokenString := newToken(userFromDB.ID, userFromDB.Role)

	responseData := map[string]string{
		"token":   tokenString,
		"message": "User logged in successfully!",
	}

	c.Data["json"] = responseData
	c.ServeJSON()

}
