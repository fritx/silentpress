package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	maxAgeMinute = 60
	ttlVendor    = 60 * maxAgeMinute
	ttlStatic    = 5 * maxAgeMinute
	ttlPostDir   = 1 * maxAgeMinute
)

var (
	_postDir      = os.Getenv("POST_DIR")
	postDirAbs, _ = filepath.Abs(_postDir)
)

func staticRoute(r *gin.Engine) {
	// Note: How to cache static files? #1222
	// https://github.com/gin-gonic/gin/issues/1222
	r.Use(func(c *gin.Context) {
		ttl := ttlStatic
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/vendor/") && path != "/vendor/blog.js" && path != "/vendor/blog.css" {
			ttl = ttlVendor
		} else if strings.HasPrefix(path, "/p/") {
			ttl = ttlPostDir
		}
		// Apply the Cache-Control header to the static files
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", ttl))
		// Continue to the next middleware or handler
		c.Next()
	})
	// Note: Inability to use '/' for static files #75
	// https://github.com/gin-gonic/gin/issues/75
	// r.Static("/", "./static")
	r.StaticFile("/favicon.ico", "./silent/blog/favicon.ico")
	r.StaticFile("/edit", "./static/edit.html")
	r.StaticFile("/admin", "./static/admin.html")
	r.StaticFile("/login", "./static/login.html")

	// r.Use(static.Serve("/", static.LocalFile("./static", false)))
	r.Use(static.Serve("/p", static.LocalFile(_postDir, false)))
	r.Use(static.Serve("/", static.LocalFile("./silent_ext", true)))
	r.Use(static.Serve("/vendor", static.LocalFile("./silent/blog/vendor", false)))
}
