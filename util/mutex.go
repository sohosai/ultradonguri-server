package util

import (
	"sync"
)

type Mutex[T any] struct {
	Mutex sync.Mutex
	Data  T
}

func NewMutex[T any](data T) *Mutex[T] {
	return &Mutex[T]{Data: data}
	// Go ではスコープ外でも参照が続く場合は自動でヒープに割り当てが発生するらしい:(
}

// 渡す handler の中で goobs.Client のロックを取るような処理を絶対に書くな
func (self *Mutex[T]) With(handler func(T) error) error {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return handler(self.Data)
}
