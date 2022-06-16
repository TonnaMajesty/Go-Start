package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const input = `
{
 "type": "sound",
 "msg": {
  "description": "dynamite",
  "authority": "the Bruce Dickinson"
 }
}
`

type Envelope struct {
	Type string
	Msg  interface{}
}

type Envelope2 struct {
	Type string
	Msg  *json.RawMessage
}

type Sound struct {
	Description string
	Authority   string
}

// 背景：通过同一通道交换多种信息的时候，我们经常需要 JSON 具有动态的，或者更合适的参数内容。
// https://zhuanlan.zhihu.com/p/136061084
func main() {
	var env Envelope
	if err := json.Unmarshal([]byte(input), &env); err != nil {
		log.Fatal(err)
	}

	// 方式一： map[string]interface{}
	// 缺点：使用一个嵌套的 map[string]interface{} 在那里乱叫会让代码变得特别丑
	// for the love of Gopher DO NOT DO THIS
	var desc string = env.Msg.(map[string]interface{})["description"].(string)
	fmt.Println(desc)


	// 方式二：json.RawMessage
	// json.RawMessage 可以让你延迟解析相应的 JSON 数据。它会将未处理的数据存储为 []byte
	// 这种方式可以让你显式控制 Msg 的解析。从而延迟到获取到 Type 的值之后，依据 Type 的值进行解析
	var env2 Envelope2
	if err := json.Unmarshal([]byte(input), &env2); err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(*env2.Msg))

	// 方式三：结合 *json.RawMessage 和 interface{}
	var msg json.RawMessage
	env3 := Envelope{
		Msg: &msg,
	}

	if err := json.Unmarshal([]byte(input), &env3); err != nil {
		log.Fatal(err)
	}
	switch env.Type {
	case "sound":
		var s Sound
		if err := json.Unmarshal(msg, &s); err != nil {
			log.Fatal(err)
		}
		var desc string = s.Description
		fmt.Println(desc)
	default:
		log.Fatalf("unknown message type: %q", env.Type)
	}
}
