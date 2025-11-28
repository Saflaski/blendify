package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

func main() {

	flag.Parse() // required
	defer glog.Flush()

	if err := godotenv.Load(".env"); err != nil {
		glog.Fatal("godotenv.Load failed")
	}

	glog.Info("Started Main")

	DB_ADDR := os.Getenv("DB_ADDR")
	if DB_ADDR == "" {
		glog.Fatal("DB_ADDR not set in env")
	}
	DB_PASS := os.Getenv("DB_PASS")
	if _, ok := os.LookupEnv("DB_PASS"); ok == false {
		glog.Fatal("DB_PASS not set in env")
	}
	DB_NUM, err := strconv.Atoi(os.Getenv("DB_NUM"))
	if err != nil {
		glog.Fatal("DB_NUM conversion to int failed", err)
	}
	DB_PROTOCOL, err := strconv.Atoi(os.Getenv("DB_PROTOCOL"))
	if err != nil {
		glog.Fatal("DB_PROTOCOL conversion to int failed", err)
	}

	DB_EXTERN_API_KEY := os.Getenv("LASTFM_API_KEY")
	if _, ok := os.LookupEnv("LASTFM_API_KEY"); ok == false {
		glog.Fatal("LASTFM_API_KEY not set in env")
	}

	JWT_EXPR_IN_SECONDS, err := strconv.Atoi(os.Getenv("JWT_EXP"))
	if err != nil {
		glog.Fatal("JWT_EXP conversion to int failed", err)
	}

	cfg := config{
		addr:                   ":3000",
		jwtExpirationInSeconds: JWT_EXPR_IN_SECONDS,
		db: dbConfig{
			addrString: DB_ADDR,
			password:   DB_PASS,
			db:         DB_NUM,
			protocol:   DB_PROTOCOL,
		},
		external: externalConfig{
			apiKey:    DB_EXTERN_API_KEY,
			lastFMURL: "https://ws.audioscrobbler.com/2.0/",
		},
	}

	api := application{
		config: cfg,
	}

	if err := api.run(api.mount()); err != nil { //Server Start
		glog.Fatal("Server failed to start.", err)
	}

}
