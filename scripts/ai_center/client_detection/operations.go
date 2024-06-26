package client_detection

import (
	"context"

	"github.com/go-courier/courier"
	"github.com/go-courier/metax"
)

type CreateAnalyze struct {
	Body CreateAnalyzeBody `in:"body" mime:"multipart"`
}

func (req CreateAnalyze) Path() string {
	return req.Body.UrlPath
}

func (CreateAnalyze) Method() string {
	return "POST"
}

func (req *CreateAnalyze) InvokeContext(ctx context.Context, c courier.Client) (*CreateAnalyzeResp, error) {
	resp := new(CreateAnalyzeResp)
	ctx = metax.ContextWithMeta(ctx, metax.MetaFromContext(ctx).With("operation", "ai.Transmit"))
	// ts := time.Now().Unix()
	_, err := c.Do(ctx, req).Into(resp)
	// requestId, ok := ctx.Value(types.RequestIDType("requestID")).(string)
	// defer func() {
	// 	if err == nil && client_ai_res_back_transmit.NeedUpload(resp.Extral) && ok && global.Config.ClientAiResBackTransmitOn {
	// 		buf, _ := json.Marshal(resp)
	// 		metaInfo := client_ai_res_back_transmit.MetaInfo{
	// 			Timestamp:      ts,
	// 			ImgId:          requestId,
	// 			ModelSrvName:   resp.ServerName,
	// 			ProjectId:      requestId,
	// 			DeviceId:       requestId,
	// 			AnalyzeResults: string(buf),
	// 		}
	// 		data_back_transmit.Send(data_back_transmit.BackTransmitData{
	// 			FileHeaders: []*multipart.FileHeader{req.Body.File},
	// 			MetaInfo:    metaInfo,
	// 		})
	// 	}

	// }()
	return resp, err
}

func (req *CreateAnalyze) Invoke(c courier.Client) (*CreateAnalyzeResp, error) {
	return req.InvokeContext(context.Background(), c)
}

type CreateExhibitAnalyze struct {
	Body CreateExhibitAnalyzeBody `in:"body" mime:"multipart"`
}

func (req CreateExhibitAnalyze) Path() string {
	return req.Body.UrlPath
}

func (CreateExhibitAnalyze) Method() string {
	return "POST"
}

func (req *CreateExhibitAnalyze) InvokeContext(ctx context.Context, c courier.Client) (*CreateAnalyzeResp, error) {
	resp := new(CreateAnalyzeResp)
	ctx = metax.ContextWithMeta(ctx, metax.MetaFromContext(ctx).With("operation", "ai.Transmit"))
	// ts := time.Now().Unix()
	_, err := c.Do(ctx, req).Into(resp)
	// requestId, ok := ctx.Value(types.RequestIDType("requestID")).(string)
	// defer func() {
	// 	if err == nil && client_ai_res_back_transmit.NeedUpload(resp.Extral) && ok && global.Config.ClientAiResBackTransmitOn {
	// 		buf, _ := json.Marshal(resp)
	// 		metaInfo := client_ai_res_back_transmit.MetaInfo{
	// 			Timestamp:      ts,
	// 			ImgId:          requestId,
	// 			ModelSrvName:   resp.ServerName,
	// 			ProjectId:      requestId,
	// 			DeviceId:       requestId,
	// 			AnalyzeResults: string(buf),
	// 		}
	// 		data_back_transmit.Send(data_back_transmit.BackTransmitData{
	// 			FileHeaders: []*multipart.FileHeader{req.Body.Standard},
	// 			MetaInfo:    metaInfo,
	// 		})
	// 	}

	// }()
	return resp, err
}

func (req *CreateExhibitAnalyze) Invoke(c courier.Client) (*CreateAnalyzeResp, error) {
	return req.InvokeContext(context.Background(), c)
}

type CreateHistoryRectangleAnalyze struct {
	Body CreateHistoryRectangleAnalyzeBody `in:"body" mime:"multipart"`
}

func (req CreateHistoryRectangleAnalyze) Path() string {
	return req.Body.UrlPath
}

func (CreateHistoryRectangleAnalyze) Method() string {
	return "POST"
}

func (req *CreateHistoryRectangleAnalyze) InvokeContext(ctx context.Context, c courier.Client) (*CreateAnalyzeResp, error) {
	resp := new(CreateAnalyzeResp)
	ctx = metax.ContextWithMeta(ctx, metax.MetaFromContext(ctx).With("operation", "ai.Transmit"))
	// ts := time.Now().Unix()
	_, err := c.Do(ctx, req).Into(resp)
	// requestId, ok := ctx.Value(types.RequestIDType("requestID")).(string)
	// defer func() {
	// 	if err == nil && client_ai_res_back_transmit.NeedUpload(resp.Extral) && ok && global.Config.ClientAiResBackTransmitOn {
	// 		buf, _ := json.Marshal(resp)
	// 		metaInfo := client_ai_res_back_transmit.MetaInfo{
	// 			Timestamp:      ts,
	// 			ImgId:          requestId,
	// 			ModelSrvName:   resp.ServerName,
	// 			ProjectId:      requestId,
	// 			DeviceId:       requestId,
	// 			AnalyzeResults: string(buf),
	// 		}
	// 		data_back_transmit.Send(data_back_transmit.BackTransmitData{
	// 			FileHeaders: []*multipart.FileHeader{req.Body.File},
	// 			MetaInfo:    metaInfo,
	// 		})
	// 	}

	// }()
	return resp, err
}

func (req *CreateHistoryRectangleAnalyze) Invoke(c courier.Client) (*CreateAnalyzeResp, error) {
	return req.InvokeContext(context.Background(), c)
}
