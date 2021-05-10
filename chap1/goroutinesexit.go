package chap1

import (
	"context"
	"fmt"
	"time"
)

func f1(exit <-chan bool) {
	fmt.Println("f1 working...")
	// 如果没有case就绪，且⚠️没有⚠️default语句，那么就会阻塞在select这里。
	select {
	case <-exit:
		fmt.Println("f1 exit")
		return
	}
}

func f2(exit <-chan bool) {
	fmt.Println("f2 working...")
	for {
		select {
		case <-exit:
			fmt.Println("f2 exit")
			return
		default:
			fmt.Println("f2 waiting...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func f3(exit <-chan bool) {
	fmt.Println("f3 working...")
	select {
	case <-exit:
		fmt.Println("f3 exit")
		return
	}
}

func f4(exit <-chan bool) {
	fmt.Println("f4 working...")
	for {
		// 如果有default语句，会在没有case就绪的情况下执行default，所以需要for不停轮询，要不然就会错过就绪的case。
		select {
		case <-exit:
			fmt.Println("f4 exit")
			return
		default:
			fmt.Println("f4 waiting...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func f5(context context.Context) {
	fmt.Println("f5 working...")
	select {
	case <-context.Done():
		fmt.Println("f5 exit")
		return
	}
}

func f6(context context.Context) {
	fmt.Println("f6 working...")
	select {
	case <-context.Done():
		fmt.Println("f6 exit")
		return
	}
}

func ExitAllGoroutines() {
	exit := make(chan bool)
	// 构建一个可以触发取消操作的Context
	cancelContext, cancel := context.WithCancel(context.Background())
	go f1(exit)
	go f2(exit)
	go f3(exit)
	go f4(exit)
	go f5(cancelContext)
	go f6(cancelContext)
	time.Sleep(1000 * time.Millisecond)
	// 当使用close方法时，所有接收这个channel的Go程都会得到一个零值。
	close(exit)
	cancel()
	// 一个Context仅包含一个键值对，想要添加键值对，需要通过当前Context构建新的Context，然后设置新的key-value
	val := context.WithValue(context.Background(), "aaa", "bbb")
	fmt.Println(val.Value("aaa"))
	time.Sleep(1000 * time.Millisecond)
}
