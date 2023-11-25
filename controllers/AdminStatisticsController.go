package controllers

import (
	"fmt"

	"github.com/MrNi8mare/word-count-bee/models"
	"github.com/MrNi8mare/word-count-bee/utils"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

type AdminStatisticsController struct {
	beego.Controller
}

func (c *AdminStatisticsController) Post() {

	tokenString := c.Ctx.Input.Header("Authorization")
	_, role, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}

	if role != 1 {
		utils.CreateErrorResponse(&c.Controller, 400, "You are not authorized to perform this action.")
	}

	db := utils.ConnectDB()
	defer db.Close()
	var fileStats FileStats
	err = c.ParseForm(&fileStats)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid form data.")
	}

	rowCount, averageTime, err := getAdminStatistics(db, fileStats.Filename)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 404, err.Error())
		return
	}

	responseData := map[string]interface{}{
		"rowCount":    rowCount,
		"averageTime": averageTime,
	}

	c.Data["json"] = responseData
	c.ServeJSON()
}

func getAdminStatistics(db *gorm.DB, filename string) (int, float64, error) {

	var processes []models.Process
	err := db.Where("file_name = ?", filename).Find(&processes).Error
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
		return 0, 0, fmt.Errorf("no processes found for the specified file")
	}

	averageTime := totalTime / float64(rowCount)

	return rowCount, averageTime, nil
}
