package controllers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
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

	if loginData.Username == "" || loginData.Password == "" {
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

	db := utils.ConnectDB()
	defer db.Close()
	var userFromDB models.User
	result := db.Where("username = ?", loginData.Username).First(&userFromDB)
	if result.Error != nil {
		c.Ctx.Output.SetStatus(401)
		errorMessage := map[string]interface{}{
			"message": "Invalid credentials. Please check your username and password.",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(loginData.Password))
	if err != nil {
		c.Ctx.Output.SetStatus(401)
		errorMessage := map[string]interface{}{
			"message": "Invalid credentials. Please check your username and password.",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
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
		"username": username,
		"time":     time.Now().Unix(),
		"exp":      expirationTime,
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		panic(err)
	}
	return tokenString
}
