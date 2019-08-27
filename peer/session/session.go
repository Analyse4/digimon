package session

import (
	"sync"
)

type Session struct {
	data sync.Map
	Conn Connection
}

func New(conn Connection) *Session {
	sess := new(Session)
	sess.data = sync.Map{}
	sess.Conn = conn
	return sess
}

func (sess *Session) Set(k string, v interface{}) {
	sess.data.Store(k, v)
}

// Get return the value stored in session for a key. or nil if no
// value is present.
func (sess *Session) Get(k string) interface{} {
	v, _ := sess.data.Load(k)
	return v
}
