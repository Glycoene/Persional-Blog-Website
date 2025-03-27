package main

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

type Userinfo struct {
	gorm.Model
	Username string `form:"username"`
	Password string `form:"password"`
}

type BlogTemplate struct {
	gorm.Model
	Author string `form:"author"`
	Title  string `form:"title"`
	Text   string `form:"text"`
}

var DB *gorm.DB

func LoginAccount(isOn bool) gin.HandlerFunc {
	var userInfo Userinfo
	return func(ctx *gin.Context) {
		haveAccount := ctx.PostForm("HaveAccount")
		err := ctx.ShouldBind(&userInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
			ctx.Abort()
			return
		}
		if !isOn || haveAccount == "false" {
			ctx.Next()
		} else {
			var userInfo_db Userinfo
			result := DB.Where("Username = ?", userInfo.Username).First(&userInfo_db)
			if result.RowsAffected == 0 {
				ctx.HTML(http.StatusOK, "login.html", gin.H{
					"Status": "用户名不存在",
				})
				ctx.Abort()
			} else if userInfo.Password != userInfo_db.Password {
				ctx.HTML(http.StatusOK, "login.html", gin.H{
					"Status": "密码错误",
				})
				ctx.Abort()
			} else {
				ctx.Next()
			}
		}
	}
}

func CreateAccount(isOn bool) gin.HandlerFunc {
	var userInfo Userinfo
	return func(ctx *gin.Context) {
		haveAccount := ctx.PostForm("HaveAccount")
		userInfo.Username = ctx.PostForm("username")
		userInfo.Password = ctx.PostForm("password")
		if !isOn {
			ctx.Next()
			return
		}
		if haveAccount == "false" {
			var userInfo_db Userinfo
			result := DB.Where("Username = ?", userInfo.Username).First(&userInfo_db)
			if result.RowsAffected == 1 {
				ctx.HTML(http.StatusOK, "register.html", gin.H{
					"Status": "用户已存在，请更改用户名或登录",
				})
				ctx.Abort()
			} else {
				DB.Create(&userInfo)
			}
		}
		ctx.Next()
	}
}

func main() {
	dsn := "root:rootpassword@tcp(127.0.0.1:3306)/Blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	DB = db
	_ = db.AutoMigrate(&Userinfo{})
	_ = db.AutoMigrate(&BlogTemplate{})

	router := gin.Default()
	router.LoadHTMLGlob("./HTML/*")

	var userInfo Userinfo

	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/login"
		router.HandleContext(ctx)
	})

	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})

	router.POST("/mainpage", LoginAccount(true), CreateAccount(true), func(ctx *gin.Context) {
		userInfo.Username = ctx.PostForm("username")
		userInfo.Password = ctx.PostForm("password")
		ctx.HTML(http.StatusOK, "mainpage.html", gin.H{
			"Username": userInfo.Username,
		})
	})

	router.POST("/addpage", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "add.html", nil)
	})

	router.POST("/addblog", func(ctx *gin.Context) {
		var blogInfo BlogTemplate
		err := ctx.ShouldBind(&blogInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
			return
		}
		blogInfo.Author = userInfo.Username
		db.Create(&blogInfo)
		ctx.HTML(http.StatusOK, "mainpage.html", gin.H{
			"Username": userInfo.Username,
		})
	})

	router.POST("/searchpage", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "search.html", nil)
	})

	router.POST("/searchblog", func(ctx *gin.Context) {
		var blogInfo BlogTemplate
		title := ctx.PostForm("title")
		result := db.Where("Title = ?", title).First(&blogInfo)
		if result.RowsAffected == 0 {
			ctx.HTML(http.StatusOK, "search.html", gin.H{
				"Status": "无结果",
			})
			return
		}
		ctx.HTML(http.StatusOK, "search.html", gin.H{
			"ID":     blogInfo.ID,
			"Title":  blogInfo.Title,
			"Text":   blogInfo.Text,
			"Author": blogInfo.Author,
		})
	})

	router.POST("/blogpage", func(ctx *gin.Context) {
		var blogInfo BlogTemplate
		err := ctx.ShouldBind(&blogInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
			return
		}
		if userInfo.Username == blogInfo.Author {
			ctx.HTML(http.StatusOK, "blogpageCU.html", gin.H{
				"ID":     blogInfo.ID,
				"Title":  blogInfo.Title,
				"Text":   blogInfo.Text,
				"Author": blogInfo.Author,
			})
		} else {
			ctx.HTML(http.StatusOK, "blogpageCtU.html", gin.H{
				"Title":  blogInfo.Title,
				"Text":   blogInfo.Text,
				"Author": blogInfo.Author,
			})
		}
	})

	router.POST("/updatepage", func(ctx *gin.Context) {
		var blogInfo BlogTemplate
		err := ctx.ShouldBind(&blogInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
			return
		}
		ctx.HTML(http.StatusOK, "update.html", gin.H{
			"ID":    blogInfo.ID,
			"Title": blogInfo.Title,
			"Text":  blogInfo.Text,
		})
	})

	router.POST("/updateblog", func(ctx *gin.Context) {
		var blogInfo BlogTemplate
		var blogInfo_db BlogTemplate
		err := ctx.ShouldBind(&blogInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
			return
		}
		result := db.First(&blogInfo_db, blogInfo.ID)
		if result.RowsAffected == 0 {
			ctx.HTML(http.StatusOK, "search.html", gin.H{
				"Status": "无结果",
			})
			return
		}
		blogInfo_db.Title = blogInfo.Title
		blogInfo_db.Text = blogInfo.Text
		db.Save(&blogInfo_db)
		ctx.HTML(http.StatusOK, "mainpage.html", gin.H{
			"Username": userInfo.Username,
		})
	})

	router.POST("/deleteblog", func(ctx *gin.Context) {
		id := ctx.PostForm("ID")
		db.Delete(&BlogTemplate{}, id)
		ctx.HTML(http.StatusOK, "mainpage.html", gin.H{
			"Username": userInfo.Username,
		})
	})

	_ = router.Run()
}
