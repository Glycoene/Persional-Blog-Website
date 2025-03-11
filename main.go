package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Userinfo struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./HTML/*")
	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", func(ctx *gin.Context) {
		var u Userinfo
		err := ctx.ShouldBind(&u)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": err.Error(),
			})
		} else if u.Password != "123" {
			ctx.HTML(http.StatusOK, "login.html", gin.H{
				"Status": "用户名或密码错误",
			})
		} else {
			ctx.HTML(http.StatusOK, "blogpage.html", gin.H{
				"Username": u.Username,
			})
		}
	})
	r.Run()
}
