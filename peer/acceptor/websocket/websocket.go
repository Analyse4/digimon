package websocket

import (
	"github.com/Analyse4/digimon/logger"
	"github.com/Analyse4/digimon/peer/session"
	"github.com/Analyse4/digimon/peer/session/wsconnection"
	"github.com/Analyse4/digimon/service"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type Websocket struct{}

var (
	log *logrus.Entry
)

func init() {
	log = logger.GetLogger().WithField("pkg", "websocket")
}

func (ws *Websocket) Accept(s service.Service) {
	urlObj, err := url.Parse(s.GetAddr())
	if err != nil {
		log.Fatalln(err)
	}

	sm, err := s.GetSessionManager()
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

		sm.Add(session.New(wsconnection.NewConnection(c)))
	})

	log.WithFields(logrus.Fields{
		"addr":    urlObj.Host,
		"service": s.GetName(),
	}).Info("service start")
	err = http.ListenAndServe(urlObj.Host, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
