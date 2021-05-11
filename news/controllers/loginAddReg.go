package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"news/models"
)

type UserController struct {
	beego.Controller
}


func (c *UserController) Get() {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "register.html"
}

func (c *UserController) Post() {
	//1.拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("pwd")
	//2.对数据进行校验
	if userName=="" || pwd == ""{
		beego.Info("数据不能为空")
		c.Redirect("/reg",302)
		return
	}
	//3.插入数据库
	o :=orm.NewOrm()

	user := models.User{}
	user.Name = userName
	user.Pwd = pwd

	_,err := o.Insert(&user) //一定是地址
	if err != nil{
		beego.Info("插入数据失败")
		c.Redirect("/reg",302)
		return
	}
	//c.Ctx.WriteString("插入成功")
	c.Redirect("/login",302)
}

/**
登录的get方法
 */
func (c *UserController) ShowLogin() {
	userName :=c.Ctx.GetCookie("userName")
	if userName != ""{
		c.Data["userName"] = userName
		c.Data["checked"] = "checked"
	}
	c.TplName = "login.html"
	//c.Ctx.SetCookie("key","value",time)
	//c.Ctx.GetCookie("key")

}

/**
登录的Post方法  业务逻辑处理
 */
func (c *UserController) HandleLogin() {
	//拿到数据
	userName := c.GetString("userName")
	remember := c.GetString("remember")
	beego.Info("remember =",remember)

	pwd := c.GetString("pwd")
	//判断数据是否合法
	if userName == "" || pwd == ""{
		beego.Info("输入的数据不合法")
		c.TplName = "login.html"
		return
	}
	//3.查询账户和密码是否正确

	o :=orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Pwd = pwd
	err :=  o.Read(&user,"Name","Pwd") //select * from name=? and Pwd=?
	if err != nil{
		beego.Info("查询失败")
		c.Redirect("/login",302)
		return
	}
	if remember == "on"{
		c.Ctx.SetCookie("userName",userName,200)
	}else {
		c.Ctx.SetCookie("userName",userName,-1)
	}
	c.SetSession("userName",userName)

	// c.Ctx.WriteString("登录成功")
	c.Redirect("/article/index",302)


}

//退出功能
func (c *UserController) LogOut(){
	c.DelSession("userName")
	c.Redirect("/login",302)
}

