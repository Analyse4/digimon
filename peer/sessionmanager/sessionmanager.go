package sessionmanager

import (
	"digimon/codec"
	"digimon/logger"
	"digimon/peer/session"
	"github.com/sirupsen/logrus"
	"sync"
)

const INVALIDID = 0

var (
	log *logrus.Entry
)

func init() {
	log = logger.GetLogger().WithField("pkg", "sessionmanager")
}

type SessionManager struct {
	mu        *sync.Mutex
	sessMap   map[int64]*session.Session
	currentID int64
	cleanConn chan int64
	codec     codec.Codec
}

func New(codecTyp string) *SessionManager {
	cd, _ := codec.Get(codecTyp)
	sm := &SessionManager{
		new(sync.Mutex),
		make(map[int64]*session.Session),
		INVALIDID,
		make(chan int64),
		cd,
	}

	log.WithFields(logrus.Fields{
		"codec":                 codecTyp,
		"current_connection_id": sm.currentID,
	}).Debug("session manager init successful")

	go func() {
		log.Debug("connection cleaner start")
		for {
			select {
			case connID := <-sm.cleanConn:
				if connID == INVALIDID {
					log.WithFields(logrus.Fields{
						"request_clean_connection_id": connID,
					}).Warn("invalid connection id")
				}
				sess := sm.sessMap[connID]
				sess.Conn.Close()

				log.WithFields(logrus.Fields{
					"connection_id": connID,
				}).Debug("connection is closed")

				delete(sm.sessMap, connID)

				log.WithFields(logrus.Fields{
					"session_id": connID,
				}).Debug("session is cleaned")
			}
		}
	}()
	return sm
}

func (sm *SessionManager) Add(sess *session.Session) {
	sm.mu.Lock()
	sm.currentID = sm.currentID + 1
	sm.sessMap[sm.currentID] = sess
	sm.mu.Unlock()
	sess.Conn.SetID(sm.currentID)
	sess.Conn.SetReqDeleteConn(sm.cleanConn)

	log.WithFields(logrus.Fields{
		"current_connection_id": sm.GetCurrentConnID(),
		"total_connection_num":  len(sm.sessMap),
	}).Debug("new connection")

	sm.connInit(sess)
}

func (sm *SessionManager) GetCurrentConnID() int64 { return sm.currentID }

func (sm *SessionManager) connInit(sess *session.Session) {
	sess.Conn.GetWaitGroup().Add(2)
	go sess.Conn.ReadLoop(sm.codec, sess)
	go sess.Conn.WriteLoop()
	go func() {
		sess.Conn.GetWaitGroup().Wait()
		if sess.Conn.GetReqDeleteConn() != nil {
			playerID := sess.Get("PLAYERID")
			sess.Conn.GetReqDeleteConn() <- sess.Conn.GetID()
		} else {
			log.WithFields(logrus.Fields{
				"current_connection_id": sess.Conn.GetID(),
				"missing_item":          "request_delete_channel",
			}).Debug("connection init error")
		}
	}()
}
