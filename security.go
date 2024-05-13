package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const (
	rpsAllReq   = 300
	rpsLoginApi = 1

	// Note: gin session: invalid memory address or nil pointer dereference? #91
	// https://github.com/gin-contrib/sessions/issues/91
	// > mind the illegal character in your session name such as ':' or '/'
	cookieName = "silentpress-sess"
)

var (
	// fix: [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
	// Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
	trustedProxies  = []string{"127.0.0.1", "::1"} // default
	_trustedProxies = os.Getenv("TRUSTED_PROXIES")

	cookiePath   = os.Getenv("COOKIE_PATH")
	cookieSecure = os.Getenv("COOKIE_SECURE") == "1"
	cookieSecret = os.Getenv("COOKIE_SECRET")

	// silent: url path with `#`,`?`, `&`, `=` - 404 Page not found
	// `%` - 400 Bad request
	regexUnsupportedPath = regexp.MustCompile(`[?#&=%]`)
	regexLikePrivatePath = regexp.MustCompile(`([\\/]|^)[.~_]`)
)

func init() {
	assertBadCookieSecret()

	if _trustedProxies != "" {
		trustedProxies = strings.Split(_trustedProxies, ",")
	}
}

func ratelimitAllReq() gin.HandlerFunc {
	return leakyBucket(rpsAllReq)
}
func ratelimitLoginApi() gin.HandlerFunc {
	return leakyBucket(rpsLoginApi)
}

func checkIllegalPathToCreate(c *gin.Context, key string) (string, bool) {
	pathAbs := filepath.Join(postDirAbs, key)
	if !isSubpathOfPostDir(pathAbs) || isUnsupportedPath(pathAbs) {
		securityLog(c, "Illegal attempt to create: key=%q, pathAbs=%q", key, pathAbs)
		return "", false
	}
	return pathAbs, true
}
func checkIllegalFileToSave(c *gin.Context, fileKey string) (string, bool) {
	fileAbs := filepath.Join(postDirAbs, fileKey)
	if !isSubpathOfPostDir(fileAbs) || isUnsupportedPath(fileAbs) {
		securityLog(c, "Illegal attempt to save: fileKey=%q, fileAbs=%q", fileKey, fileAbs)
		return "", false
	}
	return fileAbs, true
}
func checkIllegalDirToList(c *gin.Context, dirKey string) (string, bool) {
	dirAbs := filepath.Join(postDirAbs, dirKey)
	if !isUnderPostDir(dirAbs) {
		securityLog(c, "Illegal attempt to list: dirKey=%q, dirAbs=%q", dirKey, dirAbs)
		return "", false
	}
	return dirAbs, true
}

func isLikePrivatePath(path string) bool {
	return regexLikePrivatePath.MatchString(path)
}
func isUnsupportedPath(pathAbs string) bool {
	return regexUnsupportedPath.MatchString(pathAbs)
}
func isSubpathOfPostDir(pathAbs string) bool {
	return strings.HasPrefix(pathAbs, postDirAbs+string(filepath.Separator))
}
func isUnderPostDir(pathAbs string) bool {
	return strings.HasPrefix(pathAbs, postDirAbs)
}

func securityLog(c *gin.Context, layout string, params ...any) {
	if c != nil {
		prefix := "IP=" + c.ClientIP() + " | "
		params = append([]any{prefix}, params...)
	}
	log.Printf("[security] %s"+layout+"\n", params...)
}

func assertBadCookieSecret() {
	ln := 30
	if len(cookieSecret) < ln {
		log.Fatalf("Length of env.COOKIE_SECRET should be gte %d\n", ln)
	}
}
func enhancedCookieSessions() gin.HandlerFunc {
	store := cookie.NewStore([]byte(cookieSecret))
	// mind security: cookie options
	// https://github.com/gin-contrib/sessions/blob/master/session_options_go1.11.go
	store.Options(sessions.Options{
		Path:     cookiePath,
		Secure:   cookieSecure,
		HttpOnly: true,
	})
	return sessions.Sessions(cookieName, store)
}
func enhancedGinEngine() *gin.Engine {
	// r := gin.New()
	r := gin.Default() // with default middlewares
	r.SetTrustedProxies(trustedProxies)
	return r
}

// todo: plus most-used-passwords check
// https://github.com/danielmiessler/SecLists/blob/master/Passwords/2023-200_most_used_passwords.txt
func assertBadPassword() {
	ln := 10
	if len(adminPassword) < ln {
		log.Fatalf("Length of env.ADMIN_PASSWORD should be gte %d\n", ln)
	}
}
