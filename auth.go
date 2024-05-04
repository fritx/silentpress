package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

var (
	adminUsername = os.Getenv("ADMIN_USERNAME")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

func init() {
	assertBadPassword()
}

var checkAuth gin.HandlerFunc = func(c *gin.Context) {
	_, ok := getUser(c)
	if !ok {
		c.JSON(401, errorRes{"Login required"})
		c.Abort()
	}
}

func authApis(r *gin.Engine) {
	r.POST("/api/login", ratelimitLoginApi(), func(c *gin.Context) {
		json := &loginReq{}
		_ = c.BindJSON(json)
		if json.Username == adminUsername && json.Password == adminPassword {
			setUser(c, json.Username)
			c.JSON(200, successRes{})
		} else {
			c.JSON(400, errorRes{"Wrong username or password"})
			securityLog(c, "Failed attempt to login")
		}
	})
	r.GET("/api/session", checkAuth, func(c *gin.Context) {
		c.JSON(200, getConfigRes())
	})
}
