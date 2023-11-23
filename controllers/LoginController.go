package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/astaxie/beego"
	"github.com/golang-jwt/jwt"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Post() {

	var loginData models.LoginData
	err := c.ParseForm(&loginData)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	db := connectDB()
	defer db.Close()
	var userFromDB models.User
	result := db.Where("username = ?", loginData.Username).First(&userFromDB)
	if result.Error != nil {
		fmt.Println("Error querying the database:", result.Error)
		c.Ctx.Output.SetStatus(401)
		errorMessage := "Invalid credentials. Please check your username and password."
		c.Ctx.Output.Body([]byte(errorMessage))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(loginData.Password))
	if err != nil {
		log.Fatal(err)
		c.Ctx.Output.SetStatus(401)
		errorMessage := "Invalid credentials. Please check your username and password."
		c.Ctx.Output.Body([]byte(errorMessage))
		return
	}

	c.Ctx.Output.SetStatus(200)
	tokenString := newToken(loginData.Username)
	jsonData := createResponse(tokenString)
	c.Ctx.Output.Body(jsonData)

}

func createResponse(tokenString string) []byte {
	jsonData, err := json.Marshal(map[string]string{
		"token": tokenString,
	})
	if err != nil {
		panic(err)
	}
	return jsonData
}

func newToken(username string) string {
	hmacSampleSecret := []byte(utils.GoDotEnvVariable("JWT_SECRET"))
	expirationTime := time.Now().Add(3600 * time.Second).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "username",
		"time":     time.Now().Unix(),
		"exp":      expirationTime,
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func connectDB() *gorm.DB {
	dbHost := utils.GoDotEnvVariable("DB_HOST")
	dbUser := utils.GoDotEnvVariable("DB_USER")
	dbName := utils.GoDotEnvVariable("DB_NAME")
	dbPassword := utils.GoDotEnvVariable("DB_PASSWORD")
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s password=%s", dbHost, dbUser, dbName, dbPassword)

	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
