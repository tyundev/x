package socket

import (
	"runtime/debug"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type WsServer struct{}

// Write a messgae and close the connection
func (s *WsServer) WriteAndClose(ws *websocket.Conn, err error) {
	data := BuildErrorMessage("/system", err)
	ws.WriteMessage(websocket.TextMessage, data)
	ws.WriteMessage(websocket.CloseMessage, []byte{})
	ws.Close()
}

// Recover and close the connection
func (s *WsServer) Recover(ws *websocket.Conn) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			if _, ok = err.(IWebError); ok {
				s.WriteAndClose(ws, err)
				return
			}
			glog.Error(err, string(debug.Stack()))
			s.WriteAndClose(ws, errInternalServer)
		} else {
			glog.Error(r, string(debug.Stack()))
			s.WriteAndClose(ws, errInternalServer)
		}
	}
}
