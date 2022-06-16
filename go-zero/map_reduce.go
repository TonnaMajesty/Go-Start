package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/mr"
)

func main() {
	res, err := checkLegal([]int64{2,3,5,12,34,123,6,1,7})
	fmt.Println(res, err)
}

// 很多时候我们需要对一批数据进行处理，比如对一批用户id，效验每个用户的合法性并且效验过程中有一个出错就认为效验失败，返回的结果为效验合法的用户id
func checkLegal(uids []int64) ([]int64, error) {
	r, err := mr.MapReduce(func(source chan<- interface{}) {
		for _, uid := range uids {
			source <- uid
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		uid := item.(int64)
		ok, err := check(uid)
		// 如果check过程出现错误则通过cancel方法结束效验过程，并返回error整个效验过程结束
		if err != nil {
			cancel(err)
		}
		// 如果某个uid效验结果为false则最终结果不返回该uid
		if ok {
			writer.Write(uid)
		}
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		var uids []int64
		for p := range pipe {
			uids = append(uids, p.(int64))
		}
		writer.Write(uids)
	})
	if err != nil {
		log.Printf("check error: %v", err)
		return nil, err
	}

	return r.([]int64), nil
}

func check(uid int64) (bool, error) {
	// do something check user legal
	if uid > 10{
		return false, errors.New("uid is not leagal")
	}
	return true, nil
}