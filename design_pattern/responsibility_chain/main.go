package main

import "fmt"

func main() {
	cashier := &cashier{}

	//Set next for medical department
	medical := &medical{}
	medical.setNext(cashier)

	//Set next for doctor department
	doctor := &doctor{}
	doctor.setNext(medical)

	//Set next for reception department
	reception := &reception{}
	reception.setNext(doctor)

	patient := &patientV1{name: "abc"}
	//Patient visiting
	reception.execute(patient)
}

type patientV1 struct {
	name              string
	registrationDone  bool
	doctorCheckUpDone bool
	medicineDone      bool
	paymentDone       bool
}

// 处理者接口
type department interface {
	execute(*patientV1)
	setNext(department)
}

// 具体处理者
type reception struct {
	next department
}

func (r *reception) execute(p *patientV1) {
	if p.registrationDone {
		fmt.Println("Patient registration already done")
		r.next.execute(p)
		return
	}
	fmt.Println("Reception registering patientV1")
	p.registrationDone = true
	r.next.execute(p)
}

func (r *reception) setNext(next department) {
	r.next = next
}

// 具体处理者 doctor
type doctor struct {
	next department
}

func (d *doctor) execute(p *patientV1) {
	if p.doctorCheckUpDone {
		fmt.Println("Doctor checkup already done")
		d.next.execute(p)
		return
	}
	fmt.Println("Doctor checking patientV1")
	p.doctorCheckUpDone = true
	d.next.execute(p)
}

func (d *doctor) setNext(next department) {
	d.next = next
}

// 具体处理者 medical
type medical struct {
	next department
}

func (m *medical) execute(p *patientV1) {
	if p.medicineDone {
		fmt.Println("Medicine already given to patientV1")
		m.next.execute(p)
		return
	}
	fmt.Println("Medical giving medicine to patientV1")
	p.medicineDone = true
	m.next.execute(p)
}

func (m *medical) setNext(next department) {
	m.next = next
}

// 具体处理者 cashier
type cashier struct {
	next department
}

func (c *cashier) execute(p *patientV1) {
	if p.paymentDone {
		fmt.Println("Payment Done")
	}
	fmt.Println("Cashier getting money from patientV1 patientV1")
}

func (c *cashier) setNext(next department) {
	c.next = next
}
