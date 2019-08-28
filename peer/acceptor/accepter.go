package acceptor

import (
	"digimon/logger"
	"digimon/peer/acceptor/websocket"
	"digimon/service"
	"fmt"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logger.GetLogger().WithField("pkg", "acceptor")
}

type Acceptor interface {
	Accept(service.Service)
}

//TODO: Should perfect for general purpose
func Get(typ string) (Acceptor, error) {
	switch typ {
	case "ws":
		acp := new(websocket.Websocket)
		return acp, nil
	default:
		return nil, fmt.Errorf("acceptor is not registed")
	}
}
