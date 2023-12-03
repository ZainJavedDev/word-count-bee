package controllers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/gorilla/websocket"

	"github.com/astaxie/beego"
)

type UploadController struct {
	beego.Controller
}

type ResultController struct {
	beego.Controller
}

type Message struct {
	Routines int `form:"routines"`
}

var (
	resultChannel = make(chan map[string]interface{})
	mu            sync.Mutex
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
	// defer uploadedFile.Close()

	go func() {

		responseData := startProcess(header, uploadedFile, message, userID, resultChannel)
		fmt.Println("line76")

		mu.Lock()
		resultChannel <- responseData
		fmt.Println("line80")
		mu.Unlock()
	}()
	successMessage := map[string]interface{}{
		"message": "File uploaded successfully",
	}

	c.Data["json"] = successMessage
	c.ServeJSON()

}

func startProcess(header *multipart.FileHeader, uploadedFile multipart.File, message Message, userID uint, resultChannel chan map[string]interface{}) map[string]interface{} {

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
	return responseData
}

func (c *ResultController) Get() {
	fmt.Println("line135")
	mu.Lock()
	fmt.Println("line137")

	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		beego.Error("Error upgrading to WebSocket:", err)
		return
	}
	fmt.Println("line144")

	responseData := <-resultChannel
	conn.WriteJSON(responseData)
	fmt.Println("line149")

	mu.Unlock()
}
