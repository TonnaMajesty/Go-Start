package client_detection

import (
	"mime/multipart"

	"ai_center/geo"
)

type CreateAnalyzeBody struct {
	UrlPath string                `json:"-"`
	File    *multipart.FileHeader `name:"file"`
	Params  string                `name:"params"`
}

type CreateHistoryRectangleAnalyzeBody struct { //手动撸的，不是gen生成的
	UrlPath string                `json:"-"`
	File    *multipart.FileHeader `name:"file"`
	Params  string                `name:"params"`
}

type CreateExhibitAnalyzeBody struct {
	UrlPath    string                `json:"-"`
	Background *multipart.FileHeader `name:"background"`
	Standard   *multipart.FileHeader `name:"standard"`
	Detection  *multipart.FileHeader `name:"test"`
	// 检测移动距离阈值(像素，取值范围不得大于图像高宽像素最大值，默认值50)
	Distance int64 `name:"distance,omitempty" default:"50"`
	// 检测物体相似度阈值(取值范围0~1，默认值0.75)
	Similarity float64 `name:"similarity,omitempty" default:"0.75"`
	// 店铺类型
	StoreType string `name:"store_type,omitempty" default:"custom"`
}

type PositionItem struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

type Entity struct {
	Confidence  float64                `json:"confidence,omitempty"`
	EntityName  string                 `json:"label,omitempty"`
	Position    []PositionItem         `json:"position,omitempty"`
	Landmark    Landmark               `json:"landmark,omitempty"`
	PropertyMap map[string]interface{} `json:"propertyMap,omitempty"`
}

type Landmark struct {
	Landmark1 struct {
		Lx1 float64 `json:"lx1,omitempty"`
		Ly1 float64 `json:"ly1,omitempty"`
	}
	Landmark2 struct {
		Lx2 float64 `json:"lx2,omitempty"`
		Ly2 float64 `json:"ly2,omitempty"`
	}
	Landmark3 struct {
		Lx3 float64 `json:"lx3,omitempty"`
		Ly3 float64 `json:"ly3,omitempty"`
	}
	Landmark4 struct {
		Lx4 float64 `json:"lx4,omitempty"`
		Ly4 float64 `json:"ly4,omitempty"`
	}
	Landmark5 struct {
		Lx5 float64 `json:"lx5,omitempty"`
		Ly5 float64 `json:"ly5,omitempty"`
	}
}

type CreateAnalyzeResp struct {
	// 分析结果状态码
	Code int64 `json:"code"`
	// 描述信息
	// Message string `json:"message"`
	// 分析开始时间
	// StartTime interface{} `json:"start_time"`
	// 分析结束时间
	// EndTime interface{} `json:"end_time"`
	// 分析耗时（单位：秒）
	AnalyzeTime float64 `json:"inference_time"`
	// 分析流水号
	TraceID string `json:"traceID"`
	// 模型版本
	Version string `json:"version"`
	// 通信协议版本
	// ProtocolVersion string `json:"protocol_version"`
	// 通信协议版本
	// FrameworkVersion string `json:"framework_version"`
	// 模型类型
	ModelType string `json:"model_type"`
	// 扩展信息
	Extral string `json:"extral"` //里面包含了historyData
	// 服务名称
	// ServerName string `json:"server_name"`
	// 分析结果（检测模型或者聚合模型）
	AIEntityList []AIEntity `json:"entity_struct"`

	Alarm bool `json:"alarm"`
}

type AIProperty struct {
	// 属性名称
	Name string `json:"name"`
	// 属性值
	Value string `json:"value"`
	// 描述信息
	Desc string `json:"desc"`
	// 置信度
	Confidence float64 `json:"confidence"`
	// 扩展信息
	Extral string `json:"extral,omitempty"`
	// 单次分析结果
	SingleAnalyzeResult string `json:"singleAnalyzeResult"`
	// 版本
	Version string `json:"version,omitempty"`
	// 投票分数
	VoteScore float64 `json:"voteScore,omitempty"`

	// 扩展名字
	NameEx string `json:"nameEx,omitempty"`
	// 当前得分
	SingleAnalyzeScore float64 `json:"single_analyze_score,omitempty"`
	// 实体综合上报信息，决定实体是否检测到报警
	Alarm bool `json:"alarm"`
	// 是否需要上报（实体去重、多帧投票处理后的结果）
	NeedReport bool `json:"isreport"`
}

type AIEntity struct {
	// 实体名称
	Name string `json:"name"`
	// 描述信息
	Desc string `json:"desc"`
	// 实体识别位置
	Positions geo.Polygon `json:"bndbox"`
	// 置信度
	Confidence float64 `json:"confidence"`
	// 去重追踪id
	TrackID string `json:"trackID,omitempty"`
	// 扩展信息
	Extral string `json:"extral"`
	// 实体属性
	Properties []AIProperty `json:"property"`

	// 父级实体内部ID, todo, delete
	ParentIDOld int64 `json:"parent_id"`
	// 实体所属区域的ID，区域所属区域为自身
	ParentID string `json:"parent_id_new"`
	// 当前属性是否检测到报警
	Alarm bool `json:"alarm"`
	// 当前属性是否需要上报（实体去重、多帧投票处理后的结果）
	NeedReport bool `json:"isreport"`
}
