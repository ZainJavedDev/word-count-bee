package controllers

import (
	"encoding/json"
	"time"

	"github.com/MrNi8mare/word-count-bee/utils"

	"github.com/astaxie/beego"
	"github.com/golang-jwt/jwt"
)

type LoginController struct {
	beego.Controller
}

type LoginData struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (c *LoginController) Post() {

	var loginData LoginData
	err := c.ParseForm(&loginData)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	if loginData.Username == "admin" && loginData.Password == "123456" {
		c.Ctx.Output.SetStatus(200)
		tokenString := newToken()
		jsonData := createResponse(tokenString)
		c.Ctx.Output.Body(jsonData)
		return
	} else {
		c.Ctx.Output.SetStatus(401)
		errorMessage := "Invalid credentials. Please check your username and password."
		c.Ctx.Output.Body([]byte(errorMessage))
		return
	}

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

func newToken() string {
	hmacSampleSecret := []byte(utils.GoDotEnvVariable("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo":  "bar",
		"time": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		panic(err)
	}
	return tokenString
}
