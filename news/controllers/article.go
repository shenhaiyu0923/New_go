package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"time"
	"news/models"
	"fmt"
)

type ArticleController struct {
	beego.Controller
}

//显示列表功能
func (c *ArticleController) ShowIndex()  {
	/*userName := c.GetSession("userName")
	if userName == nil{
		c.Redirect("/login",302)
		return
	}*/
	//orm 查询
	o := orm.NewOrm()

	id,_ := c.GetInt("select")
	beego.Info("id=",id)
	//beego.Info("typename=",typeName)

	qs := o.QueryTable("Article")
	var articles []models.Article
	/*_,err :=qs.All(&articles)
	if err != nil{
		beego.Info("查询所有文章出错")
		return
	}*/
	//查询有多少数据
	var count int64
	count,err := qs.RelatedSel("ArticleType").Count()
	if id != 0 && id != 2{
		count,err = qs.RelatedSel("ArticleType").Filter("ArticleType__Id",id).Count()
	}

	if err != nil{
		beego.Info("查询错误")
		return
	}
	//beego.Info("count=",count)
	//每页显示多少个
	pageSize := 2
	pageCount := math.Ceil(float64(count)/float64(pageSize))

	//首页和末页
	pageIndex,err := c.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}
	start:= pageSize*(pageIndex -1)
	//1.参数pagesize 一页显示多少 2start 起始位置
	//select * from article where  ArticleType__Id = id limit
	if id == 0 || id == 2{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__Id",id).All(&articles)
	}


	//判断首页是否=1和末页是否=pagecount
	FirstPage := false
	if pageIndex == 1{
		FirstPage = true
	}

	LastPage := false
	if pageIndex == int(pageCount){
		LastPage = true
	}

	//获取类型数据
	var artiTypes []models.ArticleType
	_,err = o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("获取类型错误")
		return
	}
	c.Data["articleType"] =artiTypes
	c.Data["FirstPage"] = FirstPage
	c.Data["LastPage"] = LastPage
	c.Data["pageIndex"] = pageIndex
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["articles"] =articles
	//文章Id
	c.Data["typeid"] = id
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["contentHead"] = "head.html"
	c.Layout="layout.html"
	c.TplName = "index.html"
}

//文章添加get:showAdd;post:HandleAdd
func (c *ArticleController) ShowAdd() {
	o := orm.NewOrm()
	//获取类型数据
	var artiTypes []models.ArticleType
	_,err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("获取类型错误")
		return
	}
	c.Data["articleType"] =artiTypes
	c.Layout="layout.html"
	c.TplName = "add.html"
}

func (c *ArticleController) HandleAdd(){
	//1.拿到数据
	artiName := c.GetString("articleName")
	artiContent := c.GetString("content")

	//文件上传功能
	f,h,err := c.GetFile("uploadname")
	defer f.Close()

	//1.限定格式 png jpg
	fileext := path.Ext(h.Filename) //取出后缀
	beego.Info(fileext)
	if fileext != ".jpg" && fileext != ".png"{
		beego.Info("上传文件格式错误")
		return
	}

	//2.限制大小
	if h.Size > 40000000{
		beego.Info("上传文件过大")
		return
	}

	//3.对文件重新命名，防止重复
	filename := time.Now().Format("2006-01-02") + fileext //6-1-2 3:4:5

	if err != nil{
		beego.Info("上传文件失败")
		fmt.Println("getfile err",err)
	}else {
		c.SaveToFile("uploadname","./static/img/"+filename)
	}

	if artiName == "" || artiContent == ""{
		beego.Info("添加文章数据错误")
		return
	}

	//3.插入数据库
	o := orm.NewOrm()
	arti := models.Article{}
	arti.ArtiName = artiName
	arti.Acontent = artiContent
	arti.Aimg = "/static/img/"+filename
	//c.Ctx.WriteString("添加文章成功")

	//给文章添加类别
	id,err := c.GetInt("select")
	if err != nil{
		beego.Info("插入数据失败")
		return
	}
	artiType :=models.ArticleType{Id:id}
	o.Read(&artiType)

	arti.ArticleType = &artiType


	_,err = o.Insert(&arti)
	if err != nil{
		beego.Info("插入数据失败")
		return
	}
	//c.Ctx.WriteString("添加文章成功")
	c.Redirect("/index",302)


}

//显示内容详情页面
func (c *ArticleController) ShowContent()  {
	//0.获取当前用户
	userName := c.GetSession("userName")
	//1.获取文章ID
	id,err := c.GetInt("id")
	if err != nil{
		beego.Info("获取文章Id错误",err)
		return
	}
	//2.查询数据库对应的数据
	o:= orm.NewOrm()

	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil{
		beego.Info("查询错误",err)
		return
	}

	arti.Acount += 1
	//最近浏览
	//1、获取多对多的操作对象
	m2m := o.QueryM2M(&arti,"User")
	//2、获取要插入对象
	user := models.User{Name:userName.(string)}
	o.Read(&user,"Name")
	//3、多对多插入
	m2m.Add(user)

	o.Update(&arti)
	//o.LoadRelated(&arti,"User")
	var users []models.User
	o.QueryTable("User").Filter("Article__Article__Id",id).Distinct().All(&users)
	//3.传递数据给视图
	c.Data["users"] = users
	c.Data["article"] = arti
	c.Layout="layout.html"
	c.TplName = "content.html"

}

//显示页面编辑
func (c *ArticleController)ShowUpdate()  {
	//1.获取文章ID
	id,err :=c.GetInt("id")
	if err != nil{
		beego.Info("获取文章ID错误",err)
		return
	}
	//2.查询数据库功能
	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil{
		beego.Info("查询错误",err)
		return
	}
	//传递数据给视图
	c.Data["article"] =arti
	c.Layout="layout.html"
	c.TplName = "update.html"


}

func (c *ArticleController) HandleUpdate()  {
	id,_ := c.GetInt("id")
	artiname:=c.GetString("articleName")
	content := c.GetString("content")
	f,h,err := c.GetFile("uploadname")

	var filename string
	if err != nil{
		beego.Info("错误",err)
		c.Redirect("/article/index",302)
	}else{
		defer f.Close()

		//1.限定格式 png jpg
		fileext := path.Ext(h.Filename) //取出后缀
		beego.Info(fileext)
		if fileext != ".jpg" && fileext != ".png"{
			beego.Info("上传文件格式错误")
			return
		}

		//2.限制大小
		if h.Size > 40000000{
			beego.Info("上传文件过大")
			return
		}

		//3.对文件重新命名，防止重复
		filename = time.Now().Format("2006-01-02") + fileext //6-1-2 3:4:5
		c.SaveToFile("uploadname","./static/img/"+filename)
	}

	//对数据进行一个处理
	if artiname == "" || content == "" {
		beego.Info("更新数据获取失败")
		return
	}

	//3.更新数据
	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err  != nil{
		beego.Info("查询数据失败")
		return
	}
	arti.ArtiName = artiname
	arti.Acontent=content
	arti.Aimg = "./static/img/"+filename

	_,err = o.Update(&arti,"ArtiName","Acontent","Aimg")
	if err != nil{
		beego.Info("更新数据显示错误")
		return
	}
	c.Redirect("/article/index",302)

}

//删除操作
func (c *ArticleController) HandelDelete() {
	//1.拿到数据
	id,err := c.GetInt("id")
	if err != nil{
		beego.Info("获取id数据错误")
		return
	}
	//2.执行删除操作
	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil{
		beego.Info("查询错误")
		return
	}
	o.Delete(&arti)
	//3.返回列表页面
	c.Redirect("/article/index",302)

}

//显示添加类型页面
func (c *ArticleController)ShowAddType()  {
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	_,err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("没有获取到类型数据")
	}
	c.Data["articleType"] = artiTypes
	c.Layout="layout.html"
	c.TplName = "addType.html"
}

//显示添加类型页面
func (c *ArticleController)HandleAddType()  {
	//c.TplName = "addType.html"
	//1.获取内容
	typeName := c.GetString("typeName")
	//2.判断数据是否合法
	if typeName == ""{
		beego.Info("获取类型信息错误")
		return
	}

	//3.写入数据
	o := orm.NewOrm()
	artiType := models.ArticleType{}
	artiType.Tname = typeName
	_,err := o.Insert(&artiType)
	if err != nil{
		beego.Info("插入数据类型错误")
		return
	}
	//4.返回界面
	c.Redirect("/article/addType",302)
}
