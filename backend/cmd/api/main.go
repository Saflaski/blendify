package main

import (
	"flag"

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

	cfg := config {
		addr: ":3000",
		db: dbConfig{},
	}

	api := application{
		config: cfg,
	}

	if err:= api.run(api.mount()); err!=nil {	//Server Start
		glog.Fatal("Server failed to start.", err)
	}



}

