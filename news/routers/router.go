package routers

import (
	"news/controllers"
	"github.com/astaxie/beego"

	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.AfterExec,AfterExecFunc)
	beego.InsertFilter("/article/*",beego.BeforeRouter,beforExecFunc)
    beego.Router("/", &controllers.UserController{})
	beego.Router("/reg", &controllers.UserController{})
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
	beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAdd;post:HandleAdd")
    beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
    beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
    beego.Router("/article/delete",&controllers.ArticleController{},"get:HandelDelete")
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/logout",&controllers.UserController{},"get:LogOut")
    }


    var beforExecFunc = func(ctx *context.Context) {
    	userName := ctx.Input.Session("userName")
    	if userName == nil{
    		ctx.Redirect(302,"/login")
		}
	}

