package main

import (
	"digimon/config"
	"digimon/dao"
	"digimon/service"
	"flag"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "digimon"
	app.Usage = "a game about digimon"
	app.Author = "lzh"
	app.Version = "v0.0.0"
	app.Compiled = time.Now()

	app.Action = Start()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Start() cli.ActionFunc {
	return func(c *cli.Context) error {
		flag.Parse()
		http.Handle("/metrics", promhttp.Handler())
		glog.Info("digimon start!")
		config.Init()
		dao.Init()
		svc, _ := service.New("digimon", "protobufcdc", "ws", "ws://:2244/echo")
		svc.Start()
		return nil
	}
}
