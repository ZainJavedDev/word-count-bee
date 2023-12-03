package controllers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"

	"github.com/astaxie/beego"
)

type UploadController struct {
	beego.Controller
}

type Message struct {
	Routines int `form:"routines"`
}

func (c *UploadController) Post() {

	tokenString := c.Ctx.Input.Header("Authorization")
	if tokenString == "" {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}

	log.Default().Println("token: ")

	log.Default().Println(tokenString)
	userID, _, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}

	var message Message
	err = c.ParseForm(&message)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		return
	}

	if message.Routines <= 0 {
		utils.CreateErrorResponse(&c.Controller, 422, "Routines field is invalid")
	}

	uploadedFile, header, err := c.GetFile("file")
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "No file uploaded")
	}
	defer uploadedFile.Close()
	go startProcess(header, uploadedFile, message, userID)
	responseData := map[string]interface{}{
		"message": "File uploaded successfully",
	}

	c.Data["json"] = responseData
	c.ServeJSON()
}

func startProcess(header *multipart.FileHeader, uploadedFile multipart.File, message Message, userID uint) {

	defer uploadedFile.Close()
	uploadDir := "./uploads/"

	filePath := filepath.Join(uploadDir, header.Filename)

	outputFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, uploadedFile)
	if err != nil {
		fmt.Println(err)
	}

	totalCounts, routines, timeTaken := utils.ProcessFile(filePath, message.Routines)
	timeTakenString := timeTaken.String()

	db := utils.ConnectDB()
	defer db.Close()

	result := db.Create(&models.Process{FileName: header.Filename, Routines: message.Routines, Time: timeTaken, UserID: userID})
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	result = db.Create(&models.ProcessData{LineCount: totalCounts.LineCount, WordsCount: totalCounts.WordsCount, VowelsCount: totalCounts.VowelsCount, PunctuationCount: totalCounts.PunctuationCount, ProcessID: result.Value.(*models.Process).ID})

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	responseData := map[string]interface{}{
		"totalCounts": totalCounts,
		"routines":    routines,
		"timeTaken":   timeTakenString,
	}

	fmt.Println(responseData)
}
