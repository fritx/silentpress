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

var protectPrivate gin.HandlerFunc = func(c *gin.Context) {
	path := c.Request.URL.Path
	if isLikePrivatePath(path) && !isLoggedIn(c) {
		c.Status(404)
		c.Abort()
	}
}
var checkAuth gin.HandlerFunc = func(c *gin.Context) {
	if !isLoggedIn(c) {
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

func isLoggedIn(c *gin.Context) bool {
	_, ok := getUser(c)
	return ok
}
