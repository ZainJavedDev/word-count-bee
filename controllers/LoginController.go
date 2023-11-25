package controllers

import (
	"encoding/json"
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

	c.Ctx.Output.SetStatus(200)
	tokenString := newToken(userFromDB.ID, userFromDB.Role)
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

func newToken(userID uint, userRole int) string {
	hmacSampleSecret := []byte(utils.GoDotEnvVariable("JWT_SECRET"))
	expirationTime := time.Now().Add(3600 * time.Second).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userID,
		"role": userRole,
		"time": time.Now().Unix(),
		"exp":  expirationTime,
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		panic(err)
	}
	return tokenString
}
