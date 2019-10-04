package sessionmanager

import (
	"digimon/codec"
	"digimon/logger"
	"digimon/peer/cleaner"
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
	cleanConn chan *cleaner.CleanerMeta
	codec     codec.Codec
}

func New(codecTyp string) *SessionManager {
	cd, _ := codec.Get(codecTyp)
	sm := &SessionManager{
		new(sync.Mutex),
		make(map[int64]*session.Session),
		INVALIDID,
		nil,
		cd,
	}

	log.WithFields(logrus.Fields{
		"codec":                 codecTyp,
		"current_connection_id": sm.currentID,
	}).Debug("session manager init successful")

	//go func() {
	//	log.Debug("connection cleaner start")
	//	for {
	//		select {
	//		case connID := <-sm.cleanConn:
	//			if connID == INVALIDID {
	//				log.WithFields(logrus.Fields{
	//					"request_clean_connection_id": connID,
	//				}).Warn("invalid connection id")
	//			}
	//			sess := sm.sessMap[connID]
	//			sess.Conn.Close()
	//
	//			log.WithFields(logrus.Fields{
	//				"connection_id": connID,
	//			}).Debug("connection is closed")
	//
	//			delete(sm.sessMap, connID)
	//
	//			log.WithFields(logrus.Fields{
	//				"session_id": connID,
	//			}).Debug("session is cleaned")
	//		}
	//	}
	//}()
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
			playerID, _ := sess.Get("PLAYERID").(uint64)
			cleanerMeta := new(cleaner.CleanerMeta)
			cleanerMeta.ConnID = sess.Conn.GetID()
			cleanerMeta.PlayerID = playerID
			sess.Conn.GetReqDeleteConn() <- cleanerMeta
		} else {
			log.WithFields(logrus.Fields{
				"current_connection_id": sess.Conn.GetID(),
				"missing_item":          "request_delete_channel",
			}).Debug("connection init error")
		}
	}()
}

func (sm *SessionManager) SetCleaner(cleaner chan *cleaner.CleanerMeta) {
	sm.cleanConn = cleaner
}

func (sm *SessionManager) Get(id int64) *session.Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.sessMap[id]
}

func (sm *SessionManager) Delete(id int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessMap, id)
}
