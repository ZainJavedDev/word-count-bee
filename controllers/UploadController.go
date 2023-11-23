package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/golang-jwt/jwt"

	"github.com/astaxie/beego"
)

type UploadController struct {
	beego.Controller
}

func (c *UploadController) Post() {

	tokenString := c.Ctx.Input.Header("Authorization")
	if !validate(tokenString) {
		c.Ctx.Output.SetStatus(400)
		errorMessage := map[string]interface{}{
			"message": "Invalid or expired token.",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
		return
	}

	var message models.Message
	err := c.ParseForm(&message)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	if message.Routines <= 0 {
		c.Ctx.Output.SetStatus(400)
		errorMessage := map[string]interface{}{
			"message": "Routines field is invalid",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
		return
	}

	uploadedFile, header, err := c.GetFile("file")
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		errorMessage := map[string]interface{}{
			"message": "No file uploaded",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
		return
	}
	defer uploadedFile.Close()

	uploadDir := "./uploads/"

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
		fmt.Println(claims["username"])
		fmt.Println(claims["time"])
		fmt.Println(claims["exp"])
	} else {
		return false
	}
	return true
}
