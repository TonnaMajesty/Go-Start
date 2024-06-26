package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-courier/sqlx/v2/datatypes"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"git.innoai.tech/ai-apps/common/modelsx"
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
}

func NewSDXJClient(host, userName, password, clientID, clientSecret string) *SDXJClient {
	client := resty.New()
	client.SetBaseURL(host)

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

type DeviceInfoResp struct {
	Successful  bool             `json:"successful"`
	ResultValue DeviceInfoResult `json:"resultValue"`
}

type DeviceInfoResult struct {
	ItemCount int          `json:"itemCount"`
	Items     []DeviceInfo `json:"items"`
}

type DeviceInfo struct {
	DeviceID    string `json:"id"`           // 设备ID
	ChannelID   string `json:"channel"`      // 通道ID
	LineName    string `json:"line_Name"`    // 线路名称
	TowerName   string `json:"tower_Name"`   // 杆塔名称
	DeviceState string `json:"device_State"` // 设备状态 0：离线 1：在线
	//UpdateTime  int64  `json:"update_Time"`  // 最新同步时间 // TODO验证字段类型 string or int
}

type DeviceStateInfoResp struct {
	Successful  bool                  `json:"successful"`
	ResultValue DeviceStateInfoResult `json:"resultValue"`
}

type DeviceStateInfoResult struct {
	DeviceID       string  `json:"devId"`          // 设备ID
	RemainElectric float64 `json:"remainElectric"` // 剩余电量, 0.00 – 1.00 保留两位小数
	OnlineState    string  `json:"onlineState"`    // 在线状态 0：离线 1：在线
}

type HistoryPictureResp struct {
	Successful  bool                 `json:"successful"`
	ResultValue HistoryPictureResult `json:"resultValue"`
}

type HistoryPictureResult struct {
	Total       int       `json:"total"`
	Pages       int       `json:"pages"`
	PageNum     int       `json:"pageNum"`
	PageSize    int       `json:"pageSize"`
	StartRow    int       `json:"startRow"`
	EndRow      int       `json:"endRow"`
	PictureList []Picture `json:"list"`
}

type Picture struct {
	ID        string `json:"id"`         // 图片ID
	LineName  string `json:"line_name"`  // 线路名称
	TowerName string `json:"tower_name"` // 杆塔名称
	//PresetNum  string `json:"preset_num"`  // 预置位编号
	DeviceID   string `json:"dev_id"`      // 设备ID
	Url        string `json:"url"`         // 图片URL
	UploadTime string `json:"upload_time"` // 图片上传时间, yyyy-MM-dd HH:mm:ss
}

func (s *SDXJClient) GetToken() (string, error) {
	//if s.accessToken != "" && s.tokenCurrentTime+s.tokenExpiresIn-60 > time.Now().Unix() {
	//	return s.accessToken, nil
	//}

	if s.accessToken != "" && s.tokenCurrentTime+60 > time.Now().Unix() { // 每分钟刷新一次
		return s.accessToken, nil
	}

	tokenResp := TokenResp{}
	resp, err := s.client.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(
			fmt.Sprintf("grant_type=password&username=%s&password=%s&client_id=%s&client_secret=%s",
				s.userName, s.password, s.clientID, s.clientSecret)).
		SetResult(&tokenResp).
		Post("auth/realms/product/protocol/openid-connect/token")

	if err != nil {
		// 判断错误类型
		switch e := err.(type) {
		case *resty.ResponseError:
			// 如果是响应错误，说明服务器返回了错误响应
			logrus.Errorf("SDXJClient GetToken Response error: %s", e)

			// 获取响应体
			if e.Response != nil {
				logrus.Errorf("SDXJClient GetToken Response Error body:%s", string(e.Response.Body()))
			}

		default:
			// 其他类型的错误，例如网络错误等
			logrus.Errorf("SDXJClient GetToken Response Error:%s", e)
		}
		return "", err
	}

	if tokenResp.AccessToken == "" {
		logrus.Errorf("SDXJClient GetToken Empty:%s", resp.String())
		return "", errors.New("SDXJClient GetToken Empty")
	}

	logrus.Info("SDXJClient GetToken success resp: ", resp.String())

	s.saveToken(tokenResp.AccessToken, tokenResp.ExpiresIn)

	return tokenResp.AccessToken, nil
}

func (s *SDXJClient) saveToken(token string, expiresIn int64) {
	s.tokenLock.Lock()
	defer s.tokenLock.Unlock()
	s.accessToken = token
	s.tokenExpiresIn = expiresIn
	s.tokenCurrentTime = time.Now().Unix()
}

func (s *SDXJClient) GetDeviceInfo(code string) ([]DeviceInfo, error) {
	token, err := s.GetToken()
	if err != nil {
		return nil, err
	}

	res := DeviceInfoResp{}
	params := map[string]interface{}{
		"pageIndex": 1,
		"pageSize":  100,
		"filter":    fmt.Sprintf("code=%s", code),
	}
	paramsStr, _ := json.Marshal(params)
	resp, err := s.client.R().SetAuthToken(token).SetQueryString(fmt.Sprintf("params=%s", paramsStr)).SetResult(&res).Get("tDeviceInfo/")
	if err != nil {
		// 判断错误类型
		switch e := err.(type) {
		case *resty.ResponseError:
			// 如果是响应错误，说明服务器返回了错误响应
			logrus.Errorf("SDXJClient GetDeviceInfo Response error: %s", e)

			// 获取响应体
			if e.Response != nil {
				logrus.Errorf("SDXJClient GetDeviceInfo Response Error body:%s", string(e.Response.Body()))
			}

		default:
			// 其他类型的错误，例如网络错误等
			logrus.Errorf("SDXJClient GetDeviceInfo Response Error:%s", e)
		}
		return nil, err
	}
	logrus.Info("SDXJClient GetDeviceInfo success resp: ", resp.String())

	if !res.Successful {
		logrus.Errorf("SDXJClient GetDeviceInfo failed resp: %s", resp.String())
		return nil, errors.New("SDXJClient GetDeviceInfo failed")
	}

	//if len(res.ResultValue.Items) == 0 {
	//	//logrus.Errorf("SDXJClient GetDeviceInfo Items Empty: %s", resp.String())
	//	return nil, errors.New("SDXJClient GetDeviceInfo Items Empty")
	//}

	return res.ResultValue.Items, nil
}

func (s *SDXJClient) GetDeviceStateInfo(deviceID string) (DeviceStateInfoResult, error) {
	token, err := s.GetToken()
	if err != nil {
		return DeviceStateInfoResult{}, err
	}

	res := DeviceStateInfoResp{}
	resp, err := s.client.R().SetAuthToken(token).SetBody(map[string]interface{}{
		"devId": deviceID,
	}).SetResult(&res).Post("sdxj/device/imagedev_signal_battery_info")
	if err != nil {
		// 判断错误类型
		switch e := err.(type) {
		case *resty.ResponseError:
			// 如果是响应错误，说明服务器返回了错误响应
			logrus.Errorf("SDXJClient GetDeviceStateInfo Response error: %s", e)

			// 获取响应体
			if e.Response != nil {
				logrus.Errorf("SDXJClient GetDeviceStateInfo Response Error body:%s", string(e.Response.Body()))
			}

		default:
			// 其他类型的错误，例如网络错误等
			logrus.Errorf("SDXJClient GetDeviceStateInfo Response Error:%s", e)
		}
		return DeviceStateInfoResult{}, err
	}

	logrus.Debug("SDXJClient GetDeviceStateInfo success resp: ", resp.String())

	if !res.Successful {
		logrus.Errorf("SDXJClient GetDeviceStateInfo failed resp: %s", resp.String())
		return DeviceStateInfoResult{}, errors.New("SDXJClient GetDeviceStateInfo failed")
	}

	return res.ResultValue, nil
}

func (s *SDXJClient) GetHistoryPicture(deviceID string, startTime, endTime modelsx.Timestamp, pageNum, pageSize int) (HistoryPictureResult, error) {
	token, err := s.GetToken()
	if err != nil {
		return HistoryPictureResult{}, err
	}

	res := HistoryPictureResp{}
	resp, err := s.client.R().SetAuthToken(token).SetResult(&res).SetBody(map[string]interface{}{
		"devId":     deviceID,
		"startTime": startTime.In(datatypes.CST).Format("2006-01-02 15:04:05"),
		"endTime":   endTime.In(datatypes.CST).Format("2006-01-02 15:04:05"),
		"pageNum":   pageNum,
		"pageSize":  pageSize,
	}).Post("sdxj/picture/piclist")
	if err != nil {
		// 判断错误类型
		switch e := err.(type) {
		case *resty.ResponseError:
			// 如果是响应错误，说明服务器返回了错误响应
			logrus.Errorf("SDXJClient GetHistoryPicture Response error: %s", e)

			// 获取响应体
			if e.Response != nil {
				logrus.Errorf("SDXJClient GetHistoryPicture Response Error body:%s", string(e.Response.Body()))
			}

		default:
			// 其他类型的错误，例如网络错误等
			logrus.Errorf("SDXJClient GetHistoryPicture Response Error:%s", e)
		}
		return HistoryPictureResult{}, err
	}

	logrus.Debug("SDXJClient GetHistoryPicture success resp: ", resp.String())
	if !res.Successful {
		logrus.Errorf("SDXJClient GetHistoryPicture failed resp: %s", resp.String())
		return HistoryPictureResult{}, errors.New("SDXJClient GetHistoryPicture failed")
	}

	return res.ResultValue, nil
}
