package controllers

import (
	"encoding/json"
	"log"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type SignupController struct {
	beego.Controller
}

type SignupData struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (c *SignupController) Post() {

	var signupData SignupData
	if err := c.ParseForm(&signupData); err != nil {
		c.Ctx.Output.SetStatus(500)
		log.Fatal(err)
	}

	if signupData.Username == "" || signupData.Password == "" {
		c.Ctx.Output.SetStatus(400)
		errorMessage := map[string]interface{}{
			"message": "Username and password are required",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
		return
	}

	hashedPassword, err := HashPassword(signupData.Password)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		log.Fatal(err)
	}

	db := utils.ConnectDB()
	defer db.Close()

	result := db.Create(&models.User{Username: signupData.Username, Password: hashedPassword})
	if result.Error != nil {
		c.Ctx.Output.SetStatus(400)
		errorMessage := map[string]interface{}{
			"message": "User already exists",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
		return
	}

	responseData := map[string]interface{}{
		"message": "User created successfully!",
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	c.Ctx.Output.Body(jsonData)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
