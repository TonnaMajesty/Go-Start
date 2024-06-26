package client_ai_service

import (
	context "context"

	github_com_go_courier_courier "github.com/go-courier/courier"
)

type ClientAiService interface {
	WithContext(context.Context) ClientAiService
	Context() context.Context
	Analyze(req *Analyze, metas ...github_com_go_courier_courier.Metadata) (*AnalyzeResp, github_com_go_courier_courier.Metadata, error)
}

func NewClientAiService(c github_com_go_courier_courier.Client, ctx context.Context) *ClientAiServiceStruct {
	return &(ClientAiServiceStruct{
		Client: c,
		ctx:    ctx,
	})
}

type ClientAiServiceStruct struct {
	Client github_com_go_courier_courier.Client
	ctx    context.Context
}

func (c *ClientAiServiceStruct) WithContext(ctx context.Context) ClientAiService {
	cc := new(ClientAiServiceStruct)
	cc.Client = c.Client
	cc.ctx = ctx
	return cc
}

func (c *ClientAiServiceStruct) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *ClientAiServiceStruct) Analyze(req *Analyze, metas ...github_com_go_courier_courier.Metadata) (*AnalyzeResp, github_com_go_courier_courier.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}
