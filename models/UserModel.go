package models

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type LoginData struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type SignupData struct {
	Username string `form:"username"`
	Password string `form:"password"`
}
