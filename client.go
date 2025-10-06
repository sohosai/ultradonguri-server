package client

import (
	"errors"
	"sync"

	"github.com/andreykaipov/goobs"
)

type SharedClient struct {
	mutex  sync.Mutex
	client *goobs.Client
}

func NewSharedClient(client *goobs.Client) *SharedClient {
	return &SharedClient{client: client} // Go ではスコープ外でも参照が続く場合は自動でヒープに割り当てが発生するらしい:(
}

// 渡す handler の中で goobs.Client のロックを取るような処理を絶対に書くな
func (self *SharedClient) With(handler func(*goobs.Client) error) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.client == nil {
		return errors.New("shared client: nil")
	}

	return handler(self.client)
}
