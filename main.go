package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

const (
	// Note: gin session: invalid memory address or nil pointer dereference? #91
	// https://github.com/gin-contrib/sessions/issues/91
	// > mind the illegal character in your session name such as ':' or '/'
	cookieName = "bec-gin-sess"
)

var (
	host          = os.Getenv("HOST")
	port          = os.Getenv("PORT")
	cookieSecret  = os.Getenv("COOKIE_SECRET")
	_postDir      = os.Getenv("POST_DIR")
	postDirAbs, _ = filepath.Abs(_postDir)
)

func init() {
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}
	ln := 30
	if len(cookieSecret) < ln {
		log.Fatalf("Length of env.COOKIE_SECRET should be gte %d\n", ln)
	}
}

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte(cookieSecret))
	r.Use(sessions.Sessions(cookieName, store))

	adminApis(r)

	// Note: Inability to use '/' for static files #75
	// https://github.com/gin-gonic/gin/issues/75
	// r.Static("/", "./static")
	r.StaticFile("/edit", "./static/edit.html")
	r.StaticFile("/admin", "./static/admin.html")
	r.StaticFile("/login", "./static/login.html")
	r.Use(static.Serve("/", static.LocalFile("./static", false)))
	r.Use(static.Serve("/p", static.LocalFile(_postDir, false)))
	r.Use(static.Serve("/", static.LocalFile("./silent_ext", true)))
	r.Use(static.Serve("/vendor", static.LocalFile("./silent/blog/vendor", false)))
	r.Use(static.Serve("/favicon.ico", static.LocalFile("./silent/blog/favicon.ico", false)))

	addr := host + ":" + port
	fmt.Printf("Trying to listen at http://%s/ ...\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to run: err=%v\n", err)
	}
}
