package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	host = os.Getenv("HOST")
	port = os.Getenv("PORT")
)

func init() {
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}
}

func main() {
	r := enhancedGinEngine()
	r.Use(ratelimitAllReq())
	r.Use(enhancedCookieSessions())

	authApis(r)
	adminApis(r)

	// Below are static resources:
	staticRoute(r)

	addr := host + ":" + port
	log.Printf("Trying to listen on http://%s/ ...\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to run. err=%v\n", err)
	}
}
