package socket

import (
	"sync"
)

type WsClientByAuth map[string]*WsClient

func newWsClientByAuth() WsClientByAuth {
	return WsClientByAuth(map[string]*WsClient{})
}

func (r WsClientByAuth) empty() bool {
	return len(r) == 0
}

func (r WsClientByAuth) add(w *WsClient) {
	r[w.UID] = w
}

func (r WsClientByAuth) remove(w *WsClient) {
	delete(r, w.UID)
}

func (r WsClientByAuth) Send(payload []byte) {
	for _, w := range r {
		w.queueForSend(payload)
	}
}

func (r WsClientByAuth) Count() int {
	return len(r)
}

//WsClientManager : [idAuth]
type WsClientManager struct {
	sync.RWMutex
	clients map[string]WsClientByAuth
}

func NewWsClientManager() *WsClientManager {
	return &WsClientManager{
		clients: map[string]WsClientByAuth{},
	}
}


func (rb *WsClientManager) ForEach(cb func(*WsClient)) {
	rb.RLock()
	defer rb.RUnlock()
	for _, byAuth := range rb.clients {
		if byAuth != nil {
			for _, w := range byAuth {
				cb(w)
			}
		}
	}
}

func (rb *WsClientManager) ForEachByAuth(id string, cb func(*WsClient)) {
	rb.RLock()
	defer rb.RUnlock()
	var byAuth = rb.clients[id]
	if byAuth != nil {
		for _, w := range byAuth {
			cb(w)
		}
	}
}

func (rb *WsClientManager) Add(w *WsClient, id string) {
	rb.Lock()
	defer rb.Unlock()
	var r = rb.clients[id]
	if r == nil {
		r = newWsClientByAuth()
		rb.clients[id] = r
	}
	r.add(w)
}

func (rb *WsClientManager) AddMany(w *WsClient, wids []string) {
	rb.Lock()
	defer rb.Unlock()
	for _, id := range wids {
		var r = rb.clients[id]
		if r == nil {
			r = newWsClientByAuth()
			rb.clients[id] = r
		}
		r.add(w)
	}
}

func (rb *WsClientManager) Remove(w *WsClient, id string) {
	rb.Lock()
	defer rb.Unlock()
	var r = rb.clients[id]
	if r == nil {
		return
	}
	r.remove(w)
	if r.empty() {
		delete(rb.clients, id)
	}
}

func (rb *WsClientManager) RemoveMany(w *WsClient, wids []string) {
	rb.Lock()
	defer rb.Unlock()
	for _, id := range wids {
		var r = rb.clients[id]
		if r == nil {
			return
		}
		r.remove(w)
	}
}

func (rb *WsClientManager) SendJson(uri string, v interface{}) {
	var payload = BuildJsonMessage(uri, v)
	rb.SendRaw(payload)
}

func (rb *WsClientManager) SendRaw(payload []byte) {
	rb.RLock()
	defer rb.RUnlock()
	for _, r := range rb.clients {
		r.Send(payload)
	}
}

// SendToGroup send data to the group with auth id
// return the number of clients in the group
func (rb *WsClientManager) SendToGroup(authID string, uri string, v interface{}) int {
	rb.RLock()
	defer rb.RUnlock()
	var r = rb.clients[authID]
	if r == nil {
		return 0
	}
	var payload = BuildJsonMessage(uri, v)
	r.Send(payload)
	return len(r)
}

func (rb *WsClientManager) Destroy() {
	rb.ForEach(func(w *WsClient) {
		w.socket.Close()
	})
}

func (rb *WsClientManager) Count() int {
	rb.RLock()
	defer rb.RUnlock()
	var c = 0
	for _, a := range rb.clients {
		c += a.Count()
	}
	return c
}
