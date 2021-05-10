package chap1

import "fmt"

// 使用带有缓冲区的chan可以做到并发数控制这一简单的操作，同时缓冲区使用量/缓冲区容量可以简单的作为系统繁忙程度的评估标准

var msg1 string

var msg2 string

var done1 = make(chan bool)

var done2 = make(chan bool)

func setup1() {
	msg1 = "channel synchronization1"
	// 前面的程序是发送开始之前的部分
	done1 <- true
	// 后面的程序是发送结束之后的部分
	return
}

func setup2() {
	msg2 = "channel synchronization2"
	// 前面的程序是接收开始之前的部分
	<-done2
	// 后面的程序是接收完成之后的部分
	return
}

func print1() {
	go setup1()
	// 前面的程序是接收开始之前的部分
	<-done1
	// 后面的程序是接收完成之后的部分
	fmt.Println(msg1)
}

func print2() {
	go setup2()
	// 前面的程序是发送开始之前的部分
	done2 <- true
	// 后面的程序是发送完成之后的部分
	fmt.Println(msg2)
}

// Print
// "发送开始先于接收完成"
// "接收开始先于发送完成"
// 如果发送开始早于接收开始
// OTHER_CODE... ---------- SEND_BEGIN --(BLOCKED)-------------- SEND_END ---------- OTHER_CODE...
// OTHER_CODE... -------------------------------- RCV_BEGIN ------ RCV_END --------- OTHER_CODE...
// 如果接收开始早于发送开始
// OTHER_CODE... ---------- RCV_BEGIN ---(BLOCKED)-------------- RCV_END ----------- OTHER_CODE...
// OTHER_CODE... -------------------------------- SEND_BEGIN ----- SEND_END -------- OTHER_CODE...
// 如果它俩同时开始；那不论谁先结束，接收的开始一定早于发送的结束，发送的开始一定早于接收的结束
// 所以可以得到上述结论：⚠️ 发送开始一定早于接收结束；接收开始一定早于发送结束。
func Print() {
	print1()
	print2()
}
