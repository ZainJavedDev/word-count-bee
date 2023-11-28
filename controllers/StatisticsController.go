package controllers

import (
	"fmt"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

type StatisticsController struct {
	beego.Controller
}

type FileStats struct {
	Filename string `form:"filename"`
}

func (c *StatisticsController) Post() {

	tokenString := c.Ctx.Input.Header("Authorization")
	userID, _, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}

	db := utils.ConnectDB()
	defer db.Close()
	var fileStats FileStats
	err = c.ParseForm(&fileStats)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid form data.")
	}

	rowCount, averageTime, err := getStatistics(db, userID, fileStats.Filename)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 404, err.Error())
		return
	}

	responseData := map[string]interface{}{
		"processCount": rowCount,
		"averageTime":  averageTime,
	}

	c.Data["json"] = responseData
	c.ServeJSON()
}

func getStatistics(db *gorm.DB, userID uint, filename string) (int, float64, error) {

	var processes []models.Process
	err := db.Where("file_name = ? AND user_id = ?", filename, userID).Find(&processes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, 0, err
	}

	var totalTime float64
	var rowCount int

	for _, process := range processes {
		totalTime += process.Time.Seconds()
		rowCount++
	}

	if rowCount == 0 {
		return 0, 0, fmt.Errorf("no processes found for the specified file and user")
	}

	averageTime := totalTime / float64(rowCount)

	return rowCount, averageTime, nil
}
