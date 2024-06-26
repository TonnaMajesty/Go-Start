package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/singleflight"
)

type SDXJClient struct {
	client           *resty.Client
	host             string
	userName         string
	password         string
	clientID         string
	clientSecret     string
	accessToken      string
	tokenLock        sync.Mutex
	tokenCurrentTime int64
	tokenExpiresIn   int64
	Count            int32
	sf               singleflight.Group
}

func NewSDXJClient(host, userName, password, clientID, clientSecret string) *SDXJClient {
	client := resty.New()
	client.SetBaseURL(host)
	client.SetTimeout(5 * time.Second)

	return &SDXJClient{
		client:       client,
		host:         host,
		userName:     userName,
		password:     password,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (s *SDXJClient) GetToken() (string, error) {
	v, err, _ := s.sf.Do("GetToken", func() (interface{}, error) {
		if s.accessToken != "" && s.tokenCurrentTime+s.tokenExpiresIn-10 > time.Now().Unix() {
			return s.accessToken, nil
		}

		var entry *base.SentinelEntry
		var blockErr *base.BlockError

		for {
			entry, blockErr = sentinel.Entry(resName, sentinel.WithTrafficType(base.Outbound))
			if blockErr != nil {
				time.Sleep(800 * time.Millisecond)
				if s.accessToken != "" {
					return s.accessToken, nil
				}
				continue
			}
			break
		}

		entry.Exit()

		tokenResp := TokenResp{
			AccessToken: randStr(10),
			ExpiresIn:   11,
		}
		fmt.Printf("============== Count ============= %d, \n", atomic.AddInt32(&s.Count, 1))
		s.saveToken(tokenResp.AccessToken, tokenResp.ExpiresIn)

		return tokenResp.AccessToken, nil
	})

	if err != nil {
		return "", err
	}

	return v.(string), nil
}

func (s *SDXJClient) saveToken(token string, expiresIn int64) {
	s.tokenLock.Lock()
	defer s.tokenLock.Unlock()
	s.accessToken = token
	s.tokenExpiresIn = expiresIn
	s.tokenCurrentTime = time.Now().Unix()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
