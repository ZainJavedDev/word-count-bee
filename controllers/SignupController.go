package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type SignupController struct {
	beego.Controller
}

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique;not null"`
	Password string
}

type SignupData struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (c *SignupController) Post() {

	var signupData SignupData
	if err := c.ParseForm(&signupData); err != nil {
		log.Fatal(err)
	}

	hashedPassword, err := HashPassword(signupData.Password)
	if err != nil {
		log.Fatal(err)
	}

	dbHost := utils.GoDotEnvVariable("DB_HOST")
	dbUser := utils.GoDotEnvVariable("DB_USER")
	dbName := utils.GoDotEnvVariable("DB_NAME")
	dbPassword := utils.GoDotEnvVariable("DB_PASSWORD")
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s password=%s", dbHost, dbUser, dbName, dbPassword)

	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&User{})
	db.Create(&User{Username: signupData.Username, Password: hashedPassword})

	fmt.Println("User Created")

	responseData := map[string]interface{}{
		"Status": "OK",
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
