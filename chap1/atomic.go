package chap1

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var count1 int32

var count2 int32

var mutex sync.Mutex

var count3 int32

var count4 int32

var queue = make(chan int)

// 没有任何同步策略
func counterNoLock() {
	for i := 0; i < 10000; i++ {
		count1++
	}
}

// 使用互斥锁实现同步操作
func counterLock() {
	defer mutex.Unlock()
	mutex.Lock()
	for i := 0; i < 10000; i++ {
		count2++
	}
}

// 使用原子操作实现自增，在仅需原子操作的情况下，性能最好
func counterAtomic() {
	for i := 0; i < 10000; i++ {
		atomic.AddInt32(&count3, 1)
	}
}

// 使用Channel实现同步原语
func counterChannel() {
	<- queue
	for i := 0; i < 10000; i++ {
		count4++
	}
	queue <- 0
}

func Counter() {
	count1 = 0
	count2 = 0
	count3 = 0
	count4 = 0
	for i := 0; i < 10; i++ {
		go counterLock()
	}
	for i := 0; i < 10; i++ {
		go counterNoLock()
	}
	for i := 0; i < 10; i++ {
		go counterAtomic()
	}
	for i := 0; i < 10; i++ {
		go counterChannel()
	}
	queue <- 0
	time.Sleep(1 * time.Second)
	<- queue
	fmt.Println(count1)
	fmt.Println(count2)
	fmt.Println(count3)
	fmt.Println(count4)
}
