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
	beego.Router("/upload", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/signup", &controllers.SignupController{})
}
