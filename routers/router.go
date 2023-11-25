package routers

import (
	"github.com/MrNi8mare/word-count-bee/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/signup", &controllers.SignupController{})

	beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/processes", &controllers.ProcessController{})
	beego.Router("/statistics", &controllers.StatisticsController{})

	beego.Router("/admin/login", &controllers.LoginController{})
	beego.Router("/admin/statistics", &controllers.AdminStatisticsController{})
	beego.Router("/admin/processes", &controllers.AdminProcessController{})

	beego.Router("/migrate", &controllers.MigrationController{})
}
