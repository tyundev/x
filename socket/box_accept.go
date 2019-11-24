package socket

import (
	"bytes"
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"runtime/debug"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 8) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 64 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var (
	// DefaultUpgrader specifies the parameters for upgrading an HTTP
	// connection to a WebSocket connection.
	DefaultUpgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// DefaultDialer is a dialer with all fields set to the default zero values.
	DefaultDialer = websocket.DefaultDialer
)

func (b *Box) AcceptDefault(conn *websocket.Conn, a Auth) {
	b.Accept(conn, a, b.Join, b.Leave)
}

func (b *Box) Accept(conn *websocket.Conn, a Auth, onJoin func(*WsClient) error, onLeave func(*WsClient)) {
	// boxLog.Infof(0, "accept %s", a.ID())
	var c = newWsClient(a, conn)
	if onJoin != nil {
		if err := onJoin(c); err != nil {
			message := BuildErrorMessage("/system", err)
			conn.WriteMessage(websocket.TextMessage, message)
			c.socket.Close()
			return
		}
	}

	if onLeave != nil {
		defer onLeave(c)
	}
	// add the client
	ctx, cancel := context.WithCancel(context.Background())
	b.Clients.Add(c, c.Auth.ID())
	defer func() {
		cancel()
		b.Clients.Remove(c, c.Auth.ID())
		c.socket.Close()
	}()

	done := make(chan struct{}, 2)

	go func() {
		b.readMessage(c)
		done <- struct{}{}
	}()
	go func() {
		b.writeMessage(ctx, c)
		done <- struct{}{}
	}()
	<-done
}

func (b *Box) writeMessage(ctx context.Context, c *WsClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case data, ok := <-c.chanSend:
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.socket.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.socket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (b *Box) readMessage(c *WsClient) {
	c.socket.SetReadLimit(maxMessageSize)
	c.socket.SetReadDeadline(time.Now().Add(pongWait))
	c.socket.SetPongHandler(func(string) error {
		c.socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			// 	boxLog.Errorf("read message: %s", err.Error())
			// }
			break
		}
		data := bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		r, err := NewRequest(c, data)

		// boxLog.Infof(0, "ID: %s \tData: %s\n", b.ID, string(r.Payload))
		if err != nil {
			c.WriteError(err)
		} else {
			b.Serve(r)
		}
	}
}

var (
	errHandlerNotFound = BadRequest("HANDLER NOT FOUND")
	errInternalServer  = InternalServerError("SERVER ERROR")
)

func (b *Box) notFound(r *Request) {
	if r.isError() {
		boxLog.Errorln("Handler not found : " + r.RawURI)
	} else {
		r.ReplyError(errHandlerNotFound)
	}
}

func (b *Box) defaultRecover(r *Request, rc interface{}) {
	if err, ok := rc.(error); ok {
		if _, ok = err.(IWebError); ok {
			r.ReplyError(err)
			return
		}
		boxLog.Error(err, string(debug.Stack()))
		r.ReplyError(errInternalServer)
	} else {
		boxLog.Error(rc, string(debug.Stack()))
		r.ReplyError(errInternalServer)
	}
}

func (b *Box) join(w *WsClient) error {
	return nil
}

func (b *Box) leave(w *WsClient) {

}
