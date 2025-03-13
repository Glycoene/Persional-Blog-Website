package main

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

type Userinfo struct {
	Username string `form:"username" gorm:"primarykey"`
	Password string `form:"password"`
}

var db *gorm.DB

func isHaveAccount(isOpen bool) gin.HandlerFunc {
	var userInfo Userinfo
	return func(ctx *gin.Context) {
		err := ctx.ShouldBind(&userInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
			ctx.Abort()
			return
		}
		if !isOpen {
			ctx.Next()
		} else {
			var userInfo_db Userinfo
			result := db.Where("Username = ?", userInfo.Username).First(&userInfo_db)
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
				ctx.Set("Username", userInfo_db.Username)
				ctx.Next()
			}
		}
	}
}

func main() {
	dsn := "root:Sb143843819438@tcp(127.0.0.1:3306)/Blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db.AutoMigrate(&Userinfo{})
	router := gin.Default()
	router.LoadHTMLGlob("./HTML/*")
	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})
	router.POST("/login", isHaveAccount(true), func(ctx *gin.Context) {
		username := ctx.MustGet("Username")
		ctx.HTML(http.StatusOK, "blogpage.html", gin.H{
			"Username": username,
		})
	})
	router.Run()
}
