package main

import (
	"github.com/Analyse4/digimon/config"
	"github.com/Analyse4/digimon/dao"
	"github.com/Analyse4/digimon/handler"
	"github.com/Analyse4/digimon/prometheus"
	_ "github.com/Analyse4/digimon/svcregister"
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
		prometheus.Init()
		http.Handle("/metrics", promhttp.Handler())
		config.Init()
		dao.Init()
		svc := new(handler.Digimon)
		svc.Init("digimon", "protobuf", "ws", "ws://:2244/echo")
		svc.Start()
		return nil
	}
}
