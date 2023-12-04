package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	adminUsername = os.Getenv("ADMIN_USERNAME")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

var checkAuth gin.HandlerFunc = func(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username != nil && username != "" {
		c.Next()
	} else {
		c.JSON(401, errorRes{"Login required"})
	}
}

func adminApis(r *gin.Engine) {
	r.GET("/api/session", checkAuth, func(c *gin.Context) {
		c.JSON(200, successRes{})
	})
	r.POST("/api/login", func(c *gin.Context) {
		json := &loginReq{}
		_ = c.BindJSON(&json)
		if json.Username == adminUsername && json.Password == adminPassword {
			session := sessions.Default(c)
			session.Set("username", json.Username)
			session.Save()
			c.JSON(200, successRes{})
		} else {
			c.JSON(400, errorRes{"Wrong username or password"})
		}
	})
	r.GET("/api/list", checkAuth, func(c *gin.Context) {
		dir := c.Query("dir")
		dir = filepath.Join(postDir, dir)
		dirPath, _ := filepath.Abs(dir)
		if strings.HasPrefix(dirPath, postDirPath) {
			list, _ := os.ReadDir(dirPath)
			res := &listRes{}
			for _, v := range list {
				res.List = append(res.List, listFile{v.Name(), v.IsDir()})
			}
			sort.Slice(res.List, func(i, j int) bool {
				if res.List[i].IsDir && !res.List[j].IsDir {
					return true
				}
				return false
			})
			c.JSON(200, res)
		} else {
			c.JSON(400, "Invalid dir path to list")
		}
	})
}
