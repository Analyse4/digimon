package websocket

import (
	"digimon/codec"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

type Websocket struct {
	Addr  string
	Codec codec.Codec
}

func (ws *Websocket) Accept() {
	urlObj, err := url.Parse(ws.Addr)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc(urlObj.Path, func(writer http.ResponseWriter, request *http.Request) {
		_, err := new(websocket.Upgrader).Upgrade(writer, request, nil)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("new ws connection!")
	})

	err = http.ListenAndServe(urlObj.Host, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
