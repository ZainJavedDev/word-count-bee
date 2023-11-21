package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/golang-jwt/jwt"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type Message struct {
	Routines int `form:"routines"`
}

func (c *MainController) Post() {

	tokenString := c.Ctx.Input.Header("Authorization")
	if !validate(tokenString) {
		errorMessage := "Invalid or expired token."
		c.Ctx.Output.Body([]byte(errorMessage))
		c.Ctx.Output.SetStatus(401)
		return
	}

	var message Message
	err := c.ParseForm(&message)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	uploadedFile, header, err := c.GetFile("file")
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}
	defer uploadedFile.Close()

	uploadDir := "./uploads/"
	err = os.MkdirAll(uploadDir, 0777)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	filePath := filepath.Join(uploadDir, header.Filename)

	outputFile, err := os.Create(filePath)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, uploadedFile)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	totalCounts, routines, timeTaken := utils.ProcessFile(filePath, message.Routines)
	timeTakenString := timeTaken.String()

	responseData := map[string]interface{}{
		"totalCounts": totalCounts,
		"routines":    routines,
		"timeTaken":   timeTakenString,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	c.Ctx.Output.Body(jsonData)
}

func validate(tokenString string) bool {

	hmacSampleSecret := []byte(utils.GoDotEnvVariable("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"])
		fmt.Println(claims["time"])
		fmt.Println(claims["exp"])
	} else {
		return false
	}
	return true
}
