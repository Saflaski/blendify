package main

import (
	"backend-lastfm/auth"
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
	auth.ServerStart()
}
