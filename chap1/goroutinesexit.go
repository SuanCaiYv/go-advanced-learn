package chap1

import (
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

func ExitAllGoroutines() {
	exit := make(chan bool)
	go f1(exit)
	go f2(exit)
	go f3(exit)
	go f4(exit)
	time.Sleep(1000 * time.Millisecond)
	// 当使用close方法时，所有接收这个channel的Go程都会得到一个零值。
	close(exit)
	time.Sleep(1000 * time.Millisecond)
}
