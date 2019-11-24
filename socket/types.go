package socket

import (
	"fmt"
)

type EventHandler func(uri string, v interface{})
type IBoxHandler func(r *Request)

type Auth interface {
	ID() string
}

type AuthBearer string

func (a AuthBearer) ID() string {
	return string(a)
}

type AuthRandom string

func (a AuthRandom) ID() string {
	return fmt.Sprintf("r%p", &a)
}

type ReadWriter interface {
	Read() ([]byte, bool)
	Write([]byte)
}
