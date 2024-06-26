package client_ai_service

import (
	context "context"
	mime_multipart "mime/multipart"

	github_com_go_courier_courier "github.com/go-courier/courier"
	github_com_go_courier_metax "github.com/go-courier/metax"
)

type Analyze struct {
	ModelName string `in:"path" name:"modelName"`
	Product   string `in:"path" name:"product"`
	Body      struct {
		Debug     bool                       `name:"debug,omitempty"`
		File      *mime_multipart.FileHeader `name:"file"`
		Params    string                     `name:"params,omitempty"`
		RequestID string                     `name:"requestId"`
	} `in:"body" mime:"multipart"`
}

func (Analyze) Path() string {
	return "/ai-integration/:product/:modelName/analyze"
}

func (Analyze) Method() string {
	return "POST"
}

// @StatusErr[InternalServerError][500801001][InternalServerError]
func (req *Analyze) Do(ctx context.Context, c github_com_go_courier_courier.Client, metas ...github_com_go_courier_courier.Metadata) github_com_go_courier_courier.Result {

	ctx = github_com_go_courier_metax.ContextWith(ctx, "operationID", "ai-integration.Analyze")
	return c.Do(ctx, req, metas...)

}

func (req *Analyze) InvokeContext(ctx context.Context, c github_com_go_courier_courier.Client, metas ...github_com_go_courier_courier.Metadata) (*AnalyzeResp, github_com_go_courier_courier.Metadata, error) {
	resp := new(AnalyzeResp)

	meta, err := req.Do(ctx, c, metas...).Into(resp)

	return resp, meta, err
}

func (req *Analyze) Invoke(c github_com_go_courier_courier.Client, metas ...github_com_go_courier_courier.Metadata) (*AnalyzeResp, github_com_go_courier_courier.Metadata, error) {
	return req.InvokeContext(context.Background(), c, metas...)
}
