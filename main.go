package main

import (
	"digimon/config"
	"digimon/dao"
	"digimon/service"
	"fmt"
	"github.com/urfave/cli"
	"log"
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
		fmt.Println("digimon server start!")
		config.Init()
		dao.Init()
		svc, _ := service.New("digimon", "protobufcdc", "ws", "ws://:2244/echo")
		svc.Start()
		return nil
	}
}
