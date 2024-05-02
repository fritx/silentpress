package main

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	rpsAllReq   = 300
	rpsLoginApi = 1
)

var (
	// silent: url path with `#`,`?`, `&`, `=` - 404 Page not found
	// `%` - 400 Bad request
	regexUnsupportedPath = regexp.MustCompile(`[?#&=%]`)
)

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

// todo: plus most-used-passwords check
// https://github.com/danielmiessler/SecLists/blob/master/Passwords/2023-200_most_used_passwords.txt
func assertBadPassword() {
	ln := 10
	if len(adminPassword) < ln {
		log.Fatalf("Length of env.ADMIN_PASSWORD should be gte %d\n", ln)
	}
}
