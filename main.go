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
		err := ctx.ShouldBind(&userInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
			ctx.Abort()
			return
		}
		if haveAccount == "false" {
			DB.Create(&userInfo)
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
	db.AutoMigrate(&Userinfo{})

	router := gin.Default()
	router.LoadHTMLGlob("./HTML/*")

	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})

	router.POST("/mainpage", LoginAccount(true), CreateAccount(true), func(ctx *gin.Context) {
		var userInfo Userinfo
		err := ctx.ShouldBind(&userInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.HTML(http.StatusOK, "blogpage.html", gin.H{
			"Username": userInfo.Username,
		})
	})

	router.Run()
}
