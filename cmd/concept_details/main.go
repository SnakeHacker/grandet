package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/SnakeHacker/grandet/server"
)

func main() {
	configPath := flag.String("c", "conf.yaml", "config file")

	flag.Parse()
	flag.Set("logtostderr", "true")

	conf, err := server.LoadConf(*configPath)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Info("Start new server...")
	s, err := server.New(conf)
	if err != nil {
		glog.Error(err)
		return
	}

	err = s.ResetTables(&server.ConceptDetail{})
	if err != nil {
		glog.Error(err)
		return
	}
}
