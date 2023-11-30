package controllers

import (
	"time"

	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
	"github.com/golang-jwt/jwt"
)

type RefreshTokenController struct {
	beego.Controller
}

func (c *RefreshTokenController) Post() {
	tokenString := c.Ctx.Input.Header("Authorization")
	userID, role, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}
	accessToken := newAccessToken(userID, int(role))

	responseData := map[string]string{
		"access token": accessToken,
		"message":      "Token refreshed successfully!",
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = responseData
	c.ServeJSON()
}

func newAccessToken(userID uint, userRole int) string {
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
