package sessionmanager

import (
	"digimon/codec"
	"digimon/peer/session"
	"fmt"
	"github.com/golang/glog"
	"log"
	"strconv"
	"sync"
)

const INVALIDID = 0

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
	go func() {
		log.Println("conn cleaner start!")
		for {
			select {
			case connID := <-sm.cleanConn:
				if connID == INVALIDID {
					log.Println("Invalid connID")
				}
				sess := sm.sessMap[connID]
				sess.Conn.Close()
				delete(sm.sessMap, connID)
				log.Printf("conn %d is deleted", connID)
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
	fmt.Println("Total connection num: " + strconv.Itoa(len(sm.sessMap)))

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
			sess.Conn.GetReqDeleteConn() <- sess.Conn.GetID()
		} else {
			glog.Error("connection haven't init successful!")
		}
		glog.Info("Connection %d\n is closed!", sess.Conn.GetID())
	}()
}
