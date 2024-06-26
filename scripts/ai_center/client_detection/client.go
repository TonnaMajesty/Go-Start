package client_detection

import (
	"context"
	"time"

	"git.innoai.tech/component/svcutil/confhttp"
	"github.com/go-courier/courier"
	"github.com/go-courier/httptransport/client"
)

type ClientAi interface {
	WithContext(context.Context) ClientAi
	Context() context.Context
	CreateAnalyze(req *CreateAnalyze) (*CreateAnalyzeResp, error)
	CreateHistoryRectangleAnalyze(req *CreateHistoryRectangleAnalyze) (*CreateAnalyzeResp, error)
	CreateExhibitAnalyze(req *CreateExhibitAnalyze) (*CreateAnalyzeResp, error)
}

func NewClientAi(c courier.Client, ctx context.Context) *ClientAiStruct {
	var host string
	if ep, ok := c.(*confhttp.ClientEndpoint); ok {
		host = ep.Endpoint.Host()
	}
	return &(ClientAiStruct{
		Client: c,
		ctx:    ctx,
		Host:   host,
	})
}

type ClientAiStruct struct {
	Client courier.Client
	ctx    context.Context
	Host   string
}

func (c *ClientAiStruct) WithContext(ctx context.Context) ClientAi {
	cc := new(ClientAiStruct)
	cc.Client = c.Client
	cc.ctx = ctx
	cc.Host = c.Host
	return cc
}

func (mgr *ClientAiStruct) WithHost(host string) ClientAi {
	httpClient := client.Client{
		Protocol: "http",
		Host:     host,
		Timeout:  5 * time.Second,
	}
	httpClient.SetDefaults()

	cc := new(ClientAiStruct)
	cc.Client = &httpClient
	cc.Host = host
	return cc
}

func (c *ClientAiStruct) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *ClientAiStruct) CreateAnalyze(req *CreateAnalyze) (*CreateAnalyzeResp, error) {
	// start := time.Now()
	resp, err := req.InvokeContext(c.Context(), c.Client)
	// if len(req.Body.Params) > 256 {
	// 	req.Body.Params = req.Body.Params[:256] // avoid log too long
	// }
	// if err == nil {
	// 	logrus.Info("CreateAnalyze done, req:", c.Host, ",cost:", time.Since(start), ",Params:", req.Body.Params)
	// } else {
	// 	logrus.Info("CreateAnalyze failed, Params:", req.Body.Params, ",cost:", time.Since(start), ",err:", err)
	// }
	return resp, err
}

func (c *ClientAiStruct) CreateExhibitAnalyze(req *CreateExhibitAnalyze) (*CreateAnalyzeResp, error) {
	return req.InvokeContext(c.Context(), c.Client)
}

func (c *ClientAiStruct) CreateHistoryRectangleAnalyze(req *CreateHistoryRectangleAnalyze) (*CreateAnalyzeResp, error) {
	return req.InvokeContext(c.Context(), c.Client)
}
