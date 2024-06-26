package main

import (
	"fmt"

	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
)

type Temperature struct {
	Tag   string  //标签点名称
	Value float64 //数据值
	State int64   //状态
	Event string  //报警事件
}

type Water struct {
	Tag   string //标签点名称
	Value int64  //数据值
	State int64  //状态
	Event string //报警事件
}

type Smoke struct {
	Tag   string //标签点名称
	Value int64  //数据值
	State int64  //状态
	Event string //报警事件
}

const (
	eventRule = `
rule "TemperatureRule" "温度事件计算规则"
begin
   println("/***************** 温度事件计算规则 ***************/")
   tempState = 0
   if Temperature.Value < 0{
      tempState = 1
   }else if Temperature.Value > 80{
      tempState = 2
   }
   if Temperature.State != tempState{
      if tempState == 0{
         Temperature.Event = "温度正常"
      }else if tempState == 1{
         Temperature.Event = "低温报警"
      }else{
         Temperature.Event = "高温报警"
      }
   }else{
      Temperature.Event = ""
   }
   Temperature.State = tempState
end
 
rule "WaterRule" "水浸事件计算规则"
begin
   println("/***************** 水浸事件计算规则 ***************/")
   tempState = 0
   if Water.Value != 0{
      tempState = 1
   }
   if Water.State != tempState{
      if tempState == 0{
         Water.Event = "水浸正常"
      }else{
         Water.Event = "水浸异常"
      }
   }else{
      Water.Event = ""
   }
   Water.State = tempState
end
 
rule "SmokeRule" "烟雾事件计算规则"
begin
   println("/***************** 烟雾事件计算规则 ***************/")
   tempState = 0
   if Smoke.Value != 0{
      tempState = 1
   }
   if Smoke.State != tempState{
      if tempState == 0{
         Smoke.Event = "烟雾正常"
      }else{
         Smoke.Event = "烟雾报警"
      }
   }else{
      Smoke.Event = ""
   }
   Smoke.State = tempState
end
`
)

func main() {
	temperature := &Temperature{
		Tag:   "temperature",
		Value: 90,
		State: 0,
		Event: "",
	}
	water := &Water{
		Tag:   "water",
		Value: 0,
		State: 0,
		Event: "",
	}
	smoke := &Smoke{
		Tag:   "smoke",
		Value: 1,
		State: 0,
		Event: "",
	}
	dataContext := context.NewDataContext()
	dataContext.Add("Temperature", temperature)
	dataContext.Add("Water", water)
	dataContext.Add("Smoke", smoke)
	dataContext.Add("println", fmt.Println)
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	err1 := ruleBuilder.BuildRuleFromString(eventRule)
	if err1 != nil {
		panic(err1)
	}
	eng := engine.NewGengine()
	eng.ExecuteConcurrent(ruleBuilder)
	fmt.Printf("temperature Event=%s\n", temperature.Event)
	fmt.Printf("water Event=%s\n", water.Event)
	fmt.Printf("smoke Event=%s\n", smoke.Event)
	for i := 0; i < 10; i++ {
		smoke.Value = int64(i % 3)
		eng.ExecuteConcurrent(ruleBuilder)
		fmt.Printf("smoke Event=%s\n", smoke.Event)
	}
}
