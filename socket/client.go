package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

const clientChanSendLen = 1024

type WsClient struct {
	UID      string          //UID: Socket Id
	Auth     Auth            //
	socket   *websocket.Conn //
	chanSend chan []byte
}

func newWsClient(a Auth, s *websocket.Conn) *WsClient {
	var c = &WsClient{
		Auth:     a,
		socket:   s,
		chanSend: make(chan []byte, clientChanSendLen),
	}
	c.UID = fmt.Sprintf("%p", c)
	return c
}

func (c *WsClient) queueForSend(data []byte) {
	select {
	case c.chanSend <- data:
	default:
	}
}

func (c *WsClient) WriteError(err error) {
	c.queueForSend(BuildErrorMessage("/server", err))
}

func (c *WsClient) WriteJson(uri string, v interface{}) {
	c.queueForSend(BuildJsonMessage(uri, v))
}

func (c *WsClient) WriteRaw(msg []byte) {
	c.queueForSend(msg)
}
