// Copyright 1999-2020 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
)

const resName = "example-flow-qps-resource"

func main() {
	client := NewSDXJClient("http://localhost:8080", "", "", "", "")

	// We should initialize Sentinel first.
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              1,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {

				fmt.Println(client.GetToken())
			}
		}()
	}

	// Simulate a scenario in which flow rules are updated concurrently
	//go func() {
	//	time.Sleep(time.Second * 10)
	//	_, err = flow.LoadRules([]*flow.Rule{
	//		{
	//			Resource:               resName,
	//			TokenCalculateStrategy: flow.Direct,
	//			ControlBehavior:        flow.Reject,
	//			Threshold:              80,
	//			StatIntervalInMs:       1000,
	//		},
	//	})
	//	if err != nil {
	//		log.Fatalf("Unexpected error: %+v", err)
	//		return
	//	}
	//}()
	<-ch
}
