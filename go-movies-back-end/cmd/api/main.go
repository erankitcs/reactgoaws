package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"backend/internal/storage"
	"backend/internal/storage/localstorage"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

type application struct {
	DSN             string
	Domain          string
	DB              repository.DatabaseRepo
	auth            Auth
	JWTSecret       string
	JWTIssuer       string
	JWTAudience     string
	CookieDomain    string
	APIKey          string
	rootStoragePath string
	Storage         storage.VideoStorage
	ChatManager     *ChatManager
}

func main() {
	// set application config
	var app application
	// read from command line like flag
	flag.StringVar(&app.DSN, "dns", "host=localhost port=5432 user=postgres password=postgres database=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres Connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signinig secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signinig issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signinig audience")
	flag.StringVar(&app.CookieDomain, "jwt-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.StringVar(&app.APIKey, "api-key", "xyz", "api key")
	flag.StringVar(&app.rootStoragePath, "rootstorage-path", "./moviestorage", "Root Storage Path")
	flag.Parse()
	// connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.Storage = &localstorage.LocalStorage{RootPath: app.rootStoragePath}

	//fmt.Println(app.Storage.StorageDetails())

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		//CookieName:    "__Host-refresh_token", does not work on chrom or Edge
		CookieName:   "refresh_token",
		CookieDomain: app.CookieDomain,
	}
	// Initialise Chat manager
	app.ChatManager = NewChatManager()
	// start a web server
	log.Println("Starting server on port", port)
	//http.HandleFunc("/", app.Hello)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
