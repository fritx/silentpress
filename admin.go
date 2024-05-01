package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	filePerm = 0660 // `-rw-rw----`
	dirPerm  = 0770 // `drwxrwx---`
)

func adminApis(r *gin.Engine) {
	a := r.Group("/")
	a.Use(checkAuth)
	{
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
			res.List = []listFile{}
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
				log.Printf("Failed to write file %q. err=%v\n", fileAbs, err)
				return
			}
			c.JSON(200, successRes{})
		})
		a.POST("/api/new", func(c *gin.Context) {
			key := c.GetHeader("x-wiki-file")
			if key == "" {
				c.JSON(400, errorRes{"Bad request"})
				return
			}
			// mind security
			pathAbs, ok := checkIllegalPathToCreate(key)
			if !ok {
				c.JSON(400, errorRes{"Bad request"})
				return
			}
			// 检查文件是否存在
			if _, err := os.Stat(pathAbs); os.IsNotExist(err) {
				// 文件不存在，则创建文件
				isDir := strings.HasSuffix(key, "/")
				if isDir {
					if err := os.MkdirAll(pathAbs, dirPerm); err != nil {
						c.JSON(500, errorRes{"Failed to create"})
						log.Printf("Failed to mkdir %q. err=%v\n", pathAbs, err)
						return
					}
				} else {
					filename := filepath.Base(pathAbs)
					content := fmt.Sprintf("# %s\n\n> ...", filename)
					if err := os.WriteFile(pathAbs, []byte(content), filePerm); err != nil {
						c.JSON(500, errorRes{"Failed to create"})
						log.Printf("Failed to write file %q. err=%v\n", pathAbs, err)
						return
					}
				}
				c.JSON(200, successRes{})
			} else if err != nil {
				// 其他错误（非“文件不存在”错误）
				c.JSON(500, errorRes{"Failed to create"})
				log.Printf("Failed to check if %q exists. err=%v\n", pathAbs, err)
			} else {
				// 文件已存在，不进行任何操作
				c.JSON(400, errorRes{"Path already exists"})
			}
		})
	}
}
