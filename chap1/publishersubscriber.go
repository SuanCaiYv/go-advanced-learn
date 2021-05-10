package chap1

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type (
	subscriber chan interface{}
	topic      func(v interface{}) bool
)

type Publisher struct {
	lock        sync.RWMutex
	buffer      int
	timeout     time.Duration
	subscribers map[subscriber]topic
}

// NewPublisher 构建一个发布者对象, 可以设置发布超时时间和缓存队列的长度
func newPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topic),
	}
}

// Subscribe 添加一个新的订阅者，订阅全部主题
func (p *Publisher) subscribe() chan interface{} {
	return p.subscribeTopic(nil)
}

// SubscribeTopic 添加一个新的订阅者，订阅过滤器筛选后的主题
func (p *Publisher) subscribeTopic(topic topic) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.lock.Lock()
	p.subscribers[ch] = topic
	p.lock.Unlock()
	return ch
}

// Evict 退出订阅
func (p *Publisher) exit(sub chan interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// Publish 发布一个主题
func (p *Publisher) publish(v interface{}) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// Close 关闭发布者对象，同时关闭所有的订阅者管道。
func (p *Publisher) close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

// 发送主题，可以容忍一定的超时
func (p *Publisher) sendTopic(
	sub subscriber, topic topic, v interface{}, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

func PublishAndSubscribe() {
	p := newPublisher(100*time.Millisecond, 10)
	defer p.close()

	all := p.subscribe()
	golang := p.subscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.publish("hello,  world!")
	p.publish("hello, golang!")

	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
}
