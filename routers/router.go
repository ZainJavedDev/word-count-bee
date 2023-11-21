package routers

import (
	"github.com/MrNi8mare/word-count-bee/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/upload", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
}
