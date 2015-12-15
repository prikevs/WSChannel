package main

//import (
//"fmt"
//)

type Hub struct {
	connections map[*Connection]bool
	handle      chan Message
	register    chan *Connection
	unregister  chan *Connection
}

func (h *Hub) Offline(conn *Connection) {
	if _, ok := h.connections[conn]; ok {
		if s := conn.channel; s != nil {
			s.Del(conn)
		}
		delete(h.connections, conn)
		close(conn.send)
	}
}

func (h *Hub) Online(conn *Connection) {
	h.connections[conn] = true
}

type Channel struct {
	conns map[*Connection]bool
}

func (s *Channel) Add(conn *Connection) {
	s.conns[conn] = true
	conn.channel = s
}

func (s *Channel) Del(conn *Connection) {
	delete(s.conns, conn)
}

type ChannelList struct {
	channels map[string]*Channel
}

func (sl *ChannelList) GetChannel(key string) *Channel {
	if _, ok := sl.channels[key]; !ok {
		sl.channels[key] = &Channel{conns: make(map[*Connection]bool)}
	}
	return sl.channels[key]
}

func (sl *ChannelList) GC() {
	for i := range sl.channels {
		if len(sl.channels[i].conns) == 0 {
			delete(sl.channels, i)
		}
	}
}

var channellist = ChannelList{
	channels: make(map[string]*Channel),
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
			channel := channellist.GetChannel(c.chname)
			channel.Add(c)
		case c := <-hub.unregister:
			hub.Offline(c)
			channellist.GC()
		case m := <-hub.handle:
			channel := m.conn.channel
			for c := range channel.conns {
				if c == m.conn {
					continue
				}
				select {
				case c.send <- m.msg:
				default:
					hub.Offline(c)
					channellist.GC()
				}
			}
		}
		//fmt.Println(channellist.GetChannel("test"))
	}
}
