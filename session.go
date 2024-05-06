package main

import (
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	deliveryUrl = os.Getenv("DELIVERY_URL")
)

func getConfigRes() configRes {
	return configRes{deliveryUrl}
}

func setUser(c *gin.Context, username string) {
	session := sessions.Default(c)
	session.Set("username", username)
	session.Save()
}

func getUser(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil || username == "" {
		return "", false
	}
	// mind security (progressive..)
	if username != adminUsername {
		return "", false
	}
	return username.(string), true
}
