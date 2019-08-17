package websocket

import (
	"digimon/acceptor/websocket/wsconnection"
	"digimon/codec"
	"digimon/connmanager"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

type Websocket struct {
	Addr  string
	Codec codec.Codec
}

//TODO: Only start read loop
func (ws *Websocket) Accept() {
	urlObj, err := url.Parse(ws.Addr)
	if err != nil {
		log.Fatalln(err)
	}

	cm := connmanager.New()

	http.HandleFunc(urlObj.Path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			}}).Upgrade(writer, request, nil)
		if err != nil {
			log.Fatalln(err)
		}

		cm.Add(wsconnection.NewConnection(c))
		log.Printf("new ws connection %d!", cm.GetCurrentConnID())
	})

	err = http.ListenAndServe(urlObj.Host, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
