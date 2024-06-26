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
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
)

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func main() {
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	// Register a state change listener so that we could observer the state change of the internal circuit breaker.
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, maxErrorCount=50
		{
			Resource:         ALG_TYPE__WILDFIREWORKS.String(),
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   10000,
			MinRequestAmount: 1,
			StatIntervalMs:   5000,
			Threshold:        3,
		},
		{
			Resource:         ALG_TYPE__FOREIGNMATTER.String(),
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 1,
			StatIntervalMs:   10000,
			Threshold:        3,
		},
		{
			Resource:         ALG_TYPE__AERIALDETECT.String(),
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 1,
			StatIntervalMs:   10000,
			Threshold:        3,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	logging.Info("[CircuitBreaker ErrorCount] Sentinel Go circuit breaking demo is running. You may see the pass/block metric in the metric log.")

	ctx := context.Background()
	manager := NewManger(map[AlgType]int{ALG_TYPE__WILDFIREWORKS: 4, ALG_TYPE__FOREIGNMATTER: 6, ALG_TYPE__AERIALDETECT: 8})
	manager.Run(context.Background())

	go func() {
		var i uint64 = 0
		for {
			time.Sleep(200 * time.Millisecond)
			if !manager.IsFull(ALG_TYPE__WILDFIREWORKS) {
				manager.AddPriorJob(i, ALG_TYPE__WILDFIREWORKS)
			}
			//
			if !manager.IsFull(ALG_TYPE__FOREIGNMATTER) {
				manager.AddPriorJob(i, ALG_TYPE__FOREIGNMATTER)
			}
			//
			if !manager.IsFull(ALG_TYPE__AERIALDETECT) {
				manager.AddPriorJob(i, ALG_TYPE__AERIALDETECT)
			}
			i++
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			manager.AdjustWorkerNum(ctx, ALG_TYPE__WILDFIREWORKS, rand.Intn(10))
			time.Sleep(2 * time.Second)
			manager.AdjustWorkerNum(ctx, ALG_TYPE__FOREIGNMATTER, rand.Intn(10))
			time.Sleep(5 * time.Second)
			manager.AdjustWorkerNum(ctx, ALG_TYPE__AERIALDETECT, rand.Intn(10))
		}
	}()
	select {}
}
