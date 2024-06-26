package main

import (
	"fmt"
	"time"
)

// State interface
type LightState interface {
	// 亮起当前状态的交通灯
	Light()
	// 转换到新状态的时候，调用的方法
	EnterState()
	// 设置一个状态要转变的状态
	NextLight(light *TrafficLight)
	// 检测车速
	CarPassingSpeed(*TrafficLight, int, string)
}

// Context
type TrafficLight struct {
	State      LightState
	SpeedLimit int
}

func NewSimpleTrafficLight(speedLimit int) *TrafficLight {
	return &TrafficLight{
		SpeedLimit: speedLimit,
		State:      NewRedState(),
	}
}

// 接下来我们在实现具体状态前，先定义一个DefaultLightState类型用于让具体的LightState 嵌套组合，减少公用法在每个具体 LightState 实现类中的重复实现。
// 这个技巧我们也在模版和策略模式使用过，诀窍是它只实现LightState里的通用方法的默认版，不能实现所有的方法，那样的话他也就算一个 LightState 具体实现了，
// 而这不是我们想要的，我们想要的是把接口中每个类型实现逻辑不同的关键方法以及覆盖默认版的通用方法的工作，留给具体类型去实现。
type DefaultLightState struct {
	StateName string
}

func (state *DefaultLightState) CarPassingSpeed(road *TrafficLight, speed int, licensePlate string) {
	if speed > road.SpeedLimit {
		fmt.Printf("Car with license %s was speeding\n", licensePlate)
	}
}

func (state *DefaultLightState) EnterState() {
	fmt.Println("changed state to:", state.StateName)
}

func (tl *TrafficLight) TransitionState(newState LightState) {
	tl.State = newState
	tl.State.EnterState()
}

// 接下来我们定义三个具体状态类型，去实现LightState接口

// 红灯状态
type redState struct {
	DefaultLightState
}

func NewRedState() *redState {
	state := &redState{}
	state.StateName = "RED"
	return state
}

func (state *redState) Light() {
	fmt.Println("红灯亮起，不可行驶")
}

func (state *redState) CarPassingSpeed(light *TrafficLight, speed int, licensePlate string) {
	// 红灯时不能行驶， 所以这里要重写覆盖 DefaultLightState 里定义的这个方法
	if speed > 0 {
		fmt.Printf("Car with license \"%s\" ran a red light!\n", licensePlate)
	}
}

func (state *redState) NextLight(light *TrafficLight) {
	light.TransitionState(NewGreenState())
}

// 绿灯状态
type greenState struct {
	DefaultLightState
}

func NewGreenState() *greenState {
	state := &greenState{}
	state.StateName = "GREEN"
	return state
}

func (state *greenState) Light() {
	fmt.Println("绿灯亮起，请行驶")
}

func (state *greenState) NextLight(light *TrafficLight) {
	light.TransitionState(NewAmberState())
}

// 黄灯状态
type amberState struct {
	DefaultLightState
}

func NewAmberState() *amberState {
	state := &amberState{}
	state.StateName = "AMBER"
	return state
}

func (state *amberState) Light() {
	fmt.Println("黄灯亮起，请注意")
}

func (state *amberState) NextLight(light *TrafficLight) {
	light.TransitionState(NewRedState())
}

func main() {
	trafficLight := NewSimpleTrafficLight(500)

	interval := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-interval.C:
			trafficLight.State.Light()
			trafficLight.State.CarPassingSpeed(trafficLight, 25, "CN1024")
			trafficLight.State.NextLight(trafficLight)
		default:
		}
	}
}
