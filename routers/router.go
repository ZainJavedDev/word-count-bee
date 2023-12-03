package routers

import (
	"github.com/MrNi8mare/word-count-bee/controllers"
	"github.com/astaxie/beego"
)

func init() {
	apiPrefix := "/api/v1"

	beego.Router(apiPrefix+"/login", &controllers.LoginController{})
	beego.Router(apiPrefix+"/signup", &controllers.SignupController{})
	beego.Router(apiPrefix+"/refresh", &controllers.RefreshTokenController{})

	beego.Router(apiPrefix+"/upload", &controllers.UploadController{})
	beego.Router(apiPrefix+"/result", &controllers.ResultController{})

	beego.Router(apiPrefix+"/processes", &controllers.ProcessController{})
	beego.Router(apiPrefix+"/statistics", &controllers.StatisticsController{})

	adminPrefix := apiPrefix + "/admin"

	beego.Router(adminPrefix+"/login", &controllers.AdminLoginController{})
	beego.Router(adminPrefix+"/refresh", &controllers.RefreshTokenController{})
	beego.Router(adminPrefix+"/statistics", &controllers.AdminStatisticsController{})
	beego.Router(adminPrefix+"/processes", &controllers.AdminProcessController{})

	beego.Router(apiPrefix+"/migrate", &controllers.MigrationController{})
	beego.Router(apiPrefix+"/role", &controllers.MigrationController{})
}
