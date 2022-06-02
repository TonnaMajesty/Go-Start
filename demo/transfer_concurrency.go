package main

import (
	"fmt"
	"net/http"
	"time"
)

// 并发转账，数据不安全
//解决方案 一： 加锁
//解决方案 二：channel
	//解决方案 三：select 优化

type User struct {
	Cash int
}
type Transfer struct {
	Sender    *User
	Recipient *User
	Amount    int
}

func sendCashHandler(transferchan chan Transfer) {
	var val Transfer
	for {
		val = <-transferchan
		val.Sender.sendCash(val.Recipient, val.Amount)
	}
}

func (u *User) sendCash(to *User, amount int) bool {
	if u.Cash < amount {
		return false
	}
	/* 设置延迟Sleep，当多个goroutines并行执行时,便于进行数据安全分析 */
	time.Sleep(500 * time.Millisecond)
	u.Cash = u.Cash - amount
	to.Cash = to.Cash + amount
	return true
}

// 缺点：DoS(Denial of Service服务拒绝)，如果我们的转账操作慢下来，那么不断进来的请求需要等待进行转账操作的那个协程从通道中读取新数据
// 但是这个线程忙于照顾转账操作，没有闲功夫读取通道中新数据，这个情况会导致系统容易遭受DoS攻击，外界只要发送大量请求就能让系统停止响应
func handleFunc1 (w http.ResponseWriter, r *http.Request) {
	transfer := Transfer{Sender: &me, Recipient: &you, Amount: 50}
	transferchan <- transfer
	fmt.Fprintf(w, "I have $%d\n", me.Cash)
	fmt.Fprintf(w, "You have $%d\n", you.Cash)
	fmt.Fprintf(w, "Total transferred: $%d\n", (you.Cash - 500))
}


// 这里提升了事件循环，等待不能超过10秒，等待超过timeout时间，会返回一个消息给User告诉它们请求已经接受，可能会花点时间处理，请耐心等候即可
func handlefunc2 (w http.ResponseWriter, r *http.Request) {
	transfer := Transfer{Sender: &me, Recipient: &you, Amount: 50}
	/*转账*/
	result := make(chan int)
	go func(transferchan chan<- Transfer, transfer Transfer, result chan<- int) {
		transferchan <- transfer
		result <- 1
	}(transferchan, transfer, result)

	/*用select来处理超时机制*/
	select {
	case <-result:
		fmt.Fprintf(w, "I have $%d\n", me.Cash)
		fmt.Fprintf(w, "You have $%d\n", you.Cash)
		fmt.Fprintf(w, "Total transferred: $%d\n", (you.Cash - 500))
	case <-time.After(time.Second * 10): //超时处理
		fmt.Fprintf(w, "Your request has been received, but is processing slowly")
	}
}

var me = User{Cash: 500}
var you = User{Cash: 500}
// 无缓冲通道，当你写的时候对面没读，写阻塞，读的时候对面没写，读阻塞
var transferchan = make(chan Transfer)

func main() {
	go sendCashHandler(transferchan)
	http.HandleFunc("/", handlefunc2)
	http.ListenAndServe(":8080", nil)
}


