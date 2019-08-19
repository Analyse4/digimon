package websocket

import (
	"digimon/peer/connection/wsconnection"
	"digimon/service"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

type Websocket struct{}

func (ws *Websocket) Accept(s service.Service) {
	urlObj, err := url.Parse(s.GetAddr())
	if err != nil {
		log.Fatalln(err)
	}

	cm, err := s.GetConnManager()
	if err != nil {
		log.Fatalln(err)
	}

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
