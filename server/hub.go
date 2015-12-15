package main

import (
	"fmt"
)

type Hub struct {
	connections map[*Connection]bool
	handle      chan Message
	register    chan *Connection
	unregister  chan *Connection
}

func (h *Hub) Offline(conn *Connection) {
	if _, ok := h.connections[conn]; ok {
		if s := conn.session; s != nil {
			s.Del(conn)
		}
		delete(h.connections, conn)
		close(conn.send)
	}
}

func (h *Hub) Online(conn *Connection) {
	h.connections[conn] = true
}

type Session struct {
	conns map[*Connection]bool
}

func (s *Session) Add(conn *Connection) {
	s.conns[conn] = true
	conn.session = s
}

func (s *Session) Del(conn *Connection) {
	delete(s.conns, conn)
}

type SessionList struct {
	sessions map[string]*Session
}

func (sl *SessionList) GetSession(key string) *Session {
	if _, ok := sl.sessions[key]; !ok {
		sl.sessions[key] = &Session{conns: make(map[*Connection]bool)}
	}
	return sl.sessions[key]
}

func (sl *SessionList) GC() {
	for i := range sl.sessions {
		if len(sl.sessions[i].conns) == 0 {
			delete(sl.sessions, i)
		}
	}
}

var sessionlist = SessionList{
	sessions: make(map[string]*Session),
}

var hub = Hub{
	connections: make(map[*Connection]bool),
	handle:      make(chan Message),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
}

func (h *Hub) run() {
	for {
		select {
		case c := <-hub.register:
			hub.Online(c)
			session := sessionlist.GetSession("default")
			session.Add(c)
		case c := <-hub.unregister:
			hub.Offline(c)
			sessionlist.GC()
		case m := <-hub.handle:
			session := m.conn.session
			for c := range session.conns {
				if c == m.conn {
					continue
				}
				select {
				case c.send <- m.msg:
				default:
					hub.Offline(c)
					sessionlist.GC()
				}
			}
		}
		fmt.Println(sessionlist.GetSession("default"))
	}
}
