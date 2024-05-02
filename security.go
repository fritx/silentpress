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

func checkIllegalPathToCreate(key string) (string, bool) {
	pathAbs := filepath.Join(postDirAbs, key)
	if !isSubpathOfPostDir(pathAbs) || isUnsupportedPath(pathAbs) {
		securityLog("Illegal attempt to create: key=%q, pathAbs=%q", key, pathAbs)
		return "", false
	}
	return pathAbs, true
}
func checkIllegalFileToSave(fileKey string) (string, bool) {
	fileAbs := filepath.Join(postDirAbs, fileKey)
	if !isSubpathOfPostDir(fileAbs) || isUnsupportedPath(fileAbs) {
		securityLog("Illegal attempt to save: fileKey=%q, fileAbs=%q", fileKey, fileAbs)
		return "", false
	}
	return fileAbs, true
}
func checkIllegalDirToList(dirKey string) (string, bool) {
	dirAbs := filepath.Join(postDirAbs, dirKey)
	if !isUnderPostDir(dirAbs) {
		securityLog("Illegal attempt to list: dirKey=%q, dirAbs=%q", dirKey, dirAbs)
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

func securityLog(layout string, params ...any) {
	log.Printf("[security] "+layout+"\n", params...)
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
