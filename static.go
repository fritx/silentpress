package main

import (
	"fmt"
	"log"
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
	adminOnly  = os.Getenv("ADMIN_ONLY") == "1"
	_postDir   = os.Getenv("POST_DIR")
	postDirAbs = ""
)

func init() {
	if abs, err := filepath.Abs(_postDir); err != nil {
		log.Fatalf("env.POST_DIR=%q is invalid\n", _postDir)
	} else {
		postDirAbs = abs
	}
}

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
	})
	r.Use(protectPrivate)

	// Note: Inability to use '/' for static files #75
	// https://github.com/gin-gonic/gin/issues/75
	// r.Static("/", "./static")
	r.StaticFile("/edit", "./static/edit.html")
	r.StaticFile("/admin", "./static/admin.html")
	r.StaticFile("/login", "./static/login.html")
	// r.Use(static.Serve("/", static.LocalFile("./static", false)))

	r.Use(static.Serve("/", static.LocalFile("./silent_ext", true)))
	r.Use(static.Serve("/vendor", static.LocalFile("./silent/blog/vendor", false)))
	r.StaticFile("/favicon.ico", "./silent/blog/favicon.ico")

	if adminOnly {
		r.Use(checkAuth)
	}
	r.Use(static.Serve("/p", static.LocalFile(postDirAbs, false)))
}
