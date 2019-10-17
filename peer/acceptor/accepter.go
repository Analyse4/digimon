package acceptor

import (
	"fmt"
	"github.com/Analyse4/digimon/logger"
	"github.com/Analyse4/digimon/peer/acceptor/websocket"
	"github.com/Analyse4/digimon/service"
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
