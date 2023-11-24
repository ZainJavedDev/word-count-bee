package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Process struct {
	gorm.Model
	Time     time.Duration `gorm:"not null"`
	UserID   uint          `gorm:"not null"`
	FileName string        `gorm:"not null"`
	Routines int           `gorm:"not null"`

	User User `gorm:"foreignkey:UserID"`
}

type ProcessData struct {
	gorm.Model
	LineCount        int  `gorm:"not null"`
	WordsCount       int  `gorm:"not null"`
	VowelsCount      int  `gorm:"not null"`
	PunctuationCount int  `gorm:"not null"`
	ProcessID        uint `gorm:"not null"`

	Process Process `gorm:"foreignkey:ProcessID"`
}
