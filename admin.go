package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	filePerm = 0660
)

var (
	adminUsername = os.Getenv("ADMIN_USERNAME")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

var checkAuth gin.HandlerFunc = func(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil || username == "" {
		c.Abort()
		c.JSON(401, errorRes{"Login required"})
	}
}

func adminApis(r *gin.Engine) {
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

	a := r.Group("/")
	a.Use(checkAuth)
	{
		a.GET("/api/session", func(c *gin.Context) {
			c.JSON(200, successRes{})
		})
		a.GET("/api/list", func(c *gin.Context) {
			dirKey := c.Query("dir")
			// mind security
			dirAbs, ok := checkIllegalDirToList(dirKey)
			if !ok {
				c.JSON(400, errorRes{"Bad request"})
				return
			}
			list, _ := os.ReadDir(dirAbs)
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
		})
		a.POST("/api/save", func(c *gin.Context) {
			fileKey := c.GetHeader("x-wiki-file")
			if fileKey == "" {
				c.JSON(400, errorRes{"Bad request"})
				return
			}
			// mind security
			fileAbs, ok := checkIllegalFileToSave(fileKey)
			if !ok {
				c.JSON(400, errorRes{"Bad request"})
				return
			}
			bytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(400, errorRes{"Bad request"})
				return
			}
			if err := os.WriteFile(fileAbs, bytes, filePerm); err != nil {
				c.JSON(500, errorRes{"Failed to save file"})
				log.Println("Failed to save file", err)
				return
			}
			c.JSON(200, successRes{})
		})
	}
}

func checkIllegalFileToSave(fileKey string) (string, bool) {
	fileAbs := filepath.Join(postDirAbs, fileKey)
	if !strings.HasPrefix(fileAbs, postDirAbs+string(filepath.Separator)) {
		log.Printf("Illegal attempt: fileKey=%q, fileAbs=%q\n", fileKey, fileAbs)
		return "", false
	}
	return fileAbs, true
}

func checkIllegalDirToList(dirKey string) (string, bool) {
	dirAbs := filepath.Join(postDirAbs, dirKey)
	if !strings.HasPrefix(dirAbs, postDirAbs) {
		log.Printf("Illegal attempt: dirKey=%q, dirAbs=%q\n", dirKey, dirAbs)
		return "", false
	}
	return dirAbs, true
}
