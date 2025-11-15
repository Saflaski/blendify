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
	if DB_ADDR == ""{
		glog.Fatal("DB_ADDR not set in env")
	}
	DB_PASS := os.Getenv("DB_PASS")
	if DB_PASS == ""{
		glog.Fatal("DB_PASS not set in env")
	}
	DB_NUM, err := strconv.Atoi(os.Getenv("DB_NUM"))
	if err != nil{
		glog.Fatal("DB_NUM conversion to int failed", err)
	}
	DB_PROTOCOL, err := strconv.Atoi(os.Getenv("DB_PROTOCOL"))
	if err != nil{
		glog.Fatal("DB_PROTOCOL conversion to int failed", err)
	}
	

	cfg := config {
		addr: ":3000",
		db: dbConfig{
			addrString: DB_ADDR,
			password:   DB_PASS,
			db:      	DB_NUM,
			protocol:   DB_PROTOCOL,
		},
	}

	api := application{
		config: cfg,
	}

	if err:= api.run(api.mount()); err!=nil {	//Server Start
		glog.Fatal("Server failed to start.", err)
	}



}

