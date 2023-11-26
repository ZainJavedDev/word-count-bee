package routers

import (
	"github.com/MrNi8mare/word-count-bee/controllers"
	"github.com/astaxie/beego"
)

func init() {
	apiPrefix := "/api/v1"

	beego.Router(apiPrefix+"/login", &controllers.LoginController{})
	beego.Router(apiPrefix+"/signup", &controllers.SignupController{})

	beego.Router(apiPrefix+"/upload", &controllers.UploadController{})
	beego.Router(apiPrefix+"/processes", &controllers.ProcessController{})
	beego.Router(apiPrefix+"/statistics", &controllers.StatisticsController{})

	adminPrefix := apiPrefix + "/admin"

	beego.Router(adminPrefix+"/login", &controllers.LoginController{})
	beego.Router(adminPrefix+"/statistics", &controllers.AdminStatisticsController{})
	beego.Router(adminPrefix+"/processes", &controllers.AdminProcessController{})

	beego.Router(apiPrefix+"/migrate", &controllers.MigrationController{})
}
