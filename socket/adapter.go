package socket

import (
	"encoding/json"
	"fmt"
)

type boundedQueue struct {
	len  int
	data []interface{}
}

func NewBoundedQueue(len int) *boundedQueue {
	return &boundedQueue{
		len:  len,
		data: make([]interface{}, 0),
	}
}

func (b *boundedQueue) Add(v []byte) {
	entry := map[string]interface{}{}
	err := json.Unmarshal(v, &entry)
	if err != nil {
		fmt.Printf("[bounded queue] %s", err.Error())
		return
	}
	data := append([]interface{}{entry}, b.data...)
	last := len(data)
	if last > b.len {
		last = last - 1
	}
	b.data = data[:last]
}

func (b *boundedQueue) Read(start, end int) ([]interface{}, error) {
	return b.data, nil
}

type HistoryAdapter interface {
	Read(start, end int) ([]interface{}, error)
	Add(v []byte)
}
