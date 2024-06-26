package main

import (
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/mr"
)

// https://github.com/zeromicro/go-zero/blob/master/core/mr/readme-cn.md
func main() {
	pd, err := productDetail(123, 456)
	fmt.Println(pd, err)

}

// 商品详情的结果往往会依赖用户服务、库存服务、订单服务等，为了降低依赖的耗时我们往往需要对依赖做并行处理
func productDetail(uid, pid int64) (*ProductDetail, error) {
	var pd ProductDetail
	err := mr.Finish(func() (err error) {
		pd.User, err = userRpc(uid)
		return
	}, func() (err error) {
		pd.Store, err = storeRpc(pid)
		return
	}, func() (err error) {
		pd.Order, err = orderRpc(pid)
		return
	})

	if err != nil {
		log.Printf("product detail error: %v", err)
		return nil, err
	}

	return &pd, nil
}

type ProductDetail struct {
	User  string
	Store string
	Order string
}

func userRpc(uid int64) (string, error) {
	return "Tonna", nil
}

func storeRpc(pid int64) (string, error) {
	return "TMall", nil
}

func orderRpc(pid int64) (string, error) {
	return "SN13374637343123", nil
}
