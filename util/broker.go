package util

import "sync"

type Channel struct {
	mu   sync.RWMutex
	subs map[chan []byte]struct{}
}

func NewBroker() *Channel {
	return &Channel{subs: make(map[chan []byte]struct{})}
}

func (b *Channel) Subscribe() chan []byte {
	ch := make(chan []byte, 16)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *Channel) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *Channel) Publish(msg []byte) {
	b.mu.RLock()
	for ch := range b.subs {
		select {
		case ch <- msg:
		default:
		}
	}
	b.mu.RUnlock()
}
