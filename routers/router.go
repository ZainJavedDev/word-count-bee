// add login and sign up
// add postgres db
// user can view processes but can not view others processes
// add admin who can view all users process and individually as well
// add statistics as well for a file
// users can view their file statistics but can not view other users file statistics
// admin can view all file statistics
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
	// beego.Router("/admin/processes", &controllers.AdminProcessController{})

	beego.Router("/migrate", &controllers.MigrationController{})
}

// tests:
// 1. Empty username or password in signup page
// 2. Existing username in signup page
// 3. Make sure the password is hashed in the database
// 4. Wrong username or password in login page
// 5. Generated token must have expiry time
