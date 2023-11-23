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
	userID, err := validate(tokenString)
	if err != nil {
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
	err = c.ParseForm(&message)
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

	db := utils.ConnectDB()
	defer db.Close()

	result := db.Create(&models.Process{FileName: header.Filename, Routines: message.Routines, Time: timeTaken, UserID: userID})
	if result.Error != nil {
		c.Ctx.Output.SetStatus(500)
		errorMessage := map[string]interface{}{
			"message": "Error while storing the process in the database",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
		return
	}

	result = db.Create(&models.ProcessData{LineCount: totalCounts.LineCount, WordsCount: totalCounts.WordsCount, VowelsCount: totalCounts.VowelsCount, PunctuationCount: totalCounts.PunctuationCount, ProcessID: result.Value.(*models.Process).ID})

	if result.Error != nil {
		c.Ctx.Output.SetStatus(500)
		errorMessage := map[string]interface{}{
			"message": "Error while storing the process data in the database",
		}
		jsonData, err := json.Marshal(errorMessage)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			log.Fatal(err)
		}
		c.Ctx.Output.Body(jsonData)
		return
	}

	c.Ctx.Output.Body(jsonData)
}

func validate(tokenString string) (uint, error) {

	hmacSampleSecret := []byte(utils.GoDotEnvVariable("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		fmt.Println(claims["user"])
		fmt.Println(claims["time"])
		fmt.Println(claims["exp"])

		return uint(claims["user"].(float64)), nil

	} else {
		return 0, err
	}
}
