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
		glog.Warningf("godotenv.Load failed")
	}

	glog.Info("Started Main")
	ServerStart()
}
