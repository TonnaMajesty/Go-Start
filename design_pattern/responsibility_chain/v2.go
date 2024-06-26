package main

import (
	"fmt"
)

// https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247497039&idx=1&sn=ee5ef2ca2a378e9836564da0f2eae485&chksm=fa8324d8cdf4adce942debfe07b76f656bc3963ec9a70192d0195a9e2b2df9e15589e4630438#rd
// 跟V1的区别在于将公共的SetNext和Execute方法提炼到抽象类型Next中，具体handler增加了Do方法作为实际处理逻辑，Next中的Execute调用nextHandler中的Do
func main() {
	patientHealthHandler := StartHandler{}
	//
	patient := &patient{Name: "abc"}

	// 设置病人看病的链路
	reception := Reception{}
	reception.SetNext(&Clinic{}). // 诊室看病
					SetNext(&Cashier{}). // 收费处交钱
					SetNext(&Pharmacy{}) // 药房拿药

	patientHealthHandler.SetNext(&Reception{}). // 挂号
							SetNext(&Clinic{}).  // 诊室看病
							SetNext(&Cashier{}). // 收费处交钱
							SetNext(&Pharmacy{}) // 药房拿药
	//还可以扩展，比如中间加入化验科化验，图像科拍片等等

	// 执行上面设置好的业务流程
	if err := patientHealthHandler.Execute(patient); err != nil {
		// 异常
		fmt.Println("Fail | Error:" + err.Error())
		return
	}
	// 成功
	fmt.Println("Success")
}

type PatientHandler interface {
	Execute(*patient) error
	SetNext(PatientHandler) PatientHandler
	Do(*patient) error
}

// 充当抽象类型，实现公共方法，抽象方法不实现留给实现类自己实现
type Next struct {
	nextHandler PatientHandler
}

func (n *Next) SetNext(handler PatientHandler) PatientHandler {
	n.nextHandler = handler
	return handler
}

func (n *Next) Execute(patient *patient) (err error) {
	// 调用不到外部类型的 Do 方法，所以 Next 不能实现 Do 方法
	if n.nextHandler != nil {
		// 这里执行的是nextHandler的Do，导致起始handler的Do方法执行不到，因此需要StartHandler，作为第一个Handler向下转发请求
		// 由于 Go 并不支持继承，即使Next实现了Do方法，也不能达到在父类方法中调用子类方法的效果—即在我们的例子里面用Next 类型的Execute方法调用不到外部实现类型的Do方法
		// 这是 Go 语法限制，公共方法Execute并不能像面向对象那样先调用this.Do 再调用this.nextHandler.Do
		if err = n.nextHandler.Do(patient); err != nil {
			return
		}

		return n.nextHandler.Execute(patient)
	}

	return
}

//流程中的请求类--患者
type patient struct {
	Name              string
	RegistrationDone  bool
	DoctorCheckUpDone bool
	MedicineDone      bool
	PaymentDone       bool
}

// Reception 挂号处处理器
type Reception struct {
	Next
}

func (r *Reception) Do(p *patient) (err error) {
	if p.RegistrationDone {
		fmt.Println("Patient registration already done")
		return
	}
	fmt.Println("Reception registering patientV1")
	p.RegistrationDone = true
	return
}

// Clinic 诊室处理器--用于医生给病人看病
type Clinic struct {
	Next
}

func (d *Clinic) Do(p *patient) (err error) {
	if p.DoctorCheckUpDone {
		fmt.Println("Doctor checkup already done")
		return
	}
	fmt.Println("Doctor checking patientV1")
	p.DoctorCheckUpDone = true
	return
}

// Cashier 收费处处理器
type Cashier struct {
	Next
}

func (c *Cashier) Do(p *patient) (err error) {
	if p.PaymentDone {
		fmt.Println("Payment Done")
		return
	}
	fmt.Println("Cashier getting money from patientV1 patientV1")
	p.PaymentDone = true
	return
}

// Pharmacy 药房处理器
type Pharmacy struct {
	Next
}

func (m *Pharmacy) Do(p *patient) (err error) {
	if p.MedicineDone {
		fmt.Println("Medicine already given to patientV1")
		return
	}
	fmt.Println("Pharmacy giving medicine to patientV1")
	p.MedicineDone = true
	return
}

// StartHandler 不做操作，作为第一个Handler向下转发请求
type StartHandler struct {
	Next
}

// Do 空Handler的Do
func (h *StartHandler) Do(c *patient) (err error) {
	// 空Handler 这里什么也不做 只是载体 do nothing...
	return
}
