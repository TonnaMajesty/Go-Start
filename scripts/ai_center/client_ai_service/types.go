package client_ai_service

import (
	// git_innoai_tech_ai_integration_charge_pkg_arranger_dispatch "git.innoai.tech/ai-integration/charge/pkg/arranger_dispatch"

	github_com_go_courier_statuserror "github.com/go-courier/statuserror"

	"ai_center/client_detection"
	"ai_center/geo"
)

type AnalyzeResp struct {
	AnalysisDuration string     `json:"analysisDuration"`
	Duration         string     `json:"duration"`
	Alarm            bool       `json:"alarm"`
	AIEntityList     []AIEntity `json:"entity_struct"`
	HistoryData      any        `json:"history_data"`
	// 分析扩展信息
	Extral map[string]any `json:"extral,omitempty"`
}

// type GitInnoaiTechAiIntegrationChargePkgArrangerDispatchEntity = git_innoai_tech_ai_integration_charge_pkg_arranger_dispatch.Entity

// type GitInnoaiTechAiIntegrationChargePkgArrangerDispatchProperty = git_innoai_tech_ai_integration_charge_pkg_arranger_dispatch.Property

type GithubComGoCourierStatuserrorErrorField = github_com_go_courier_statuserror.ErrorField

type GithubComGoCourierStatuserrorErrorFields = github_com_go_courier_statuserror.ErrorFields

type GithubComGoCourierStatuserrorStatusErr = github_com_go_courier_statuserror.StatusErr

func (resp *AnalyzeResp) ToOldVersion() *client_detection.CreateAnalyzeResp {
	return &client_detection.CreateAnalyzeResp{
		AIEntityList: transformEntities(resp.AIEntityList),
		Alarm:        resp.Alarm,
	}
}

func transformEntities(entityList []AIEntity) []client_detection.AIEntity {
	var out []client_detection.AIEntity
	for _, e := range entityList {
		var properties []client_detection.AIProperty
		for _, p := range e.Properties {
			properties = append(properties, client_detection.AIProperty{
				Name:               p.Name,
				NameEx:             p.NameEx,
				Value:              p.Value,
				Desc:               p.Desc,
				Confidence:         p.Confidence,
				VoteScore:          p.VoteScore,
				SingleAnalyzeScore: p.SingleAnalyzeScore,
				Extral:             p.Extral,
				Alarm:              p.Alarm,
				NeedReport:         p.NeedReport,
			})
		}
		if e.Name == "polygon" { // AI so库版本roi返回的名字变成了polygon，这里统一修改，保证后续逻辑不变
			e.Name = "area"
		}
		out = append(out, client_detection.AIEntity{
			Name:       e.Name,
			Desc:       e.Desc,
			Confidence: e.Confidence,
			TrackID:    e.TrackID,
			Positions:  e.Positions,
			Properties: properties,
			ParentID:   e.ParentID,
			Extral:     e.Extral,
			Alarm:      e.Alarm,
			NeedReport: e.NeedReport,
		})
	}
	return out
}

type AIEntity struct {
	// 实体名称
	Name string `json:"name"`
	// 实体描述信息
	Desc string `json:"desc"`
	// 实体所属区域的ID，区域所属区域为自身
	ParentID string `json:"parent_id"`
	// 实体识别位置
	Positions geo.Polygon `json:"points"`
	// 置信度
	Confidence float64 `json:"conf"`
	// 去重追踪id
	TrackID string `json:"track_id,omitempty"`
	// 实体属性
	Properties []AIProperty `json:"property"`
	// 实体综合上报信息，决定实体是否上报，true上报，false不上报
	Alarm bool `json:"alarm"`
	// 扩展信息
	Extral string `json:"extral"`
	// 是否需要上报（实体去重、多帧投票处理后的结果）
	NeedReport bool `json:"isreport"`
}

type AIProperty struct {
	// 属性名称
	Name string `json:"name"`
	// 扩展名字
	NameEx string `json:"inner_name"`
	// 属性值
	Value string `json:"value"`
	// 描述信息
	Desc string `json:"desc"`
	// 置信度
	Confidence float64 `json:"conf"`
	// 投票分数
	VoteScore float64 `json:"vote_score,omitempty"`
	// 当前得分
	SingleAnalyzeScore float64 `json:"single_analyze_score,omitempty"`
	// 当前属性是否上报，true上报，false不上报
	Alarm bool `json:"alarm"`
	// 扩展信息
	Extral string `json:"extral,omitempty"`
	// 当前属性是否需要上报（实体去重、多帧投票处理后的结果）
	NeedReport bool `json:"isreport"`
}
