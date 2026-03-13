package common

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// API Response match internal/api/apiresp/resp.go
type ApiResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// WS Request match internal/im/message_handler.go
type InboundReq struct {
	ReqIdentifier int32           `json:"req_identifier"`
	MsgIncr       string          `json:"msg_incr"`
	Data          json.RawMessage `json:"data"`
}

// WS Response match internal/im/message_handler.go
type Resp struct {
	ReqIdentifier int32           `json:"req_identifier"`
	MsgIncr       string          `json:"msg_incr"`
	Code          int             `json:"code"`
	Msg           string          `json:"msg"`
	Data          json.RawMessage `json:"data"`
}

// Consts from internal/im/constant.go
const (
	WSGetNewestSeq        = 1001
	WSPullMsgBySeqList    = 1002
	WSSendMsg             = 1003
	WSSendSignalMsg       = 1004
	WSPullMsg             = 1005
	WSGetConvMaxReadSeq   = 1006
	WsPullConvLastMessage = 1007
	WSPushMsg             = 2001
	WSKickOnlineMsg       = 2002
	WsLogoutMsg           = 2003
)

const (
	WebPlatformID = 1
)

func Register(apiAddr, username, password string) (int64, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"nickname": username,
	})
	resp, err := http.Post(apiAddr+"/user/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, err
	}
	if apiResp.Code != 0 {
		return 0, fmt.Errorf("register failed: %s (code: %d)", apiResp.Msg, apiResp.Code)
	}

	var userIDStr string
	if err := json.Unmarshal(apiResp.Data, &userIDStr); err != nil {
		return 0, err
	}
	return strconv.ParseInt(userIDStr, 10, 64)
}

func Login(apiAddr, username, password string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	resp, err := http.Post(apiAddr+"/user/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}
	if apiResp.Code != 0 {
		return "", fmt.Errorf("login failed: %s (code: %d)", apiResp.Msg, apiResp.Code)
	}

	var token string
	if err := json.Unmarshal(apiResp.Data, &token); err != nil {
		return "", err
	}
	return token, nil
}

func GetWSURL(wsAddr, token string, platformID int) string {
	u, _ := url.Parse(wsAddr)
	q := u.Query()
	q.Set("token", token)
	q.Set("platformID", strconv.Itoa(platformID))
	u.RawQuery = q.Encode()
	return u.String()
}

func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GzipDecompress(data []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return io.ReadAll(zr)
}

func GetInfo(apiAddr, token string) (int64, error) {
	req, _ := http.NewRequest("GET", apiAddr+"/user/info", nil)
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, err
	}
	if apiResp.Code != 0 {
		return 0, fmt.Errorf("get info failed: %s", apiResp.Msg)
	}

	var user struct {
		UserID int64 `json:"user_id,string"`
	}
	if err := json.Unmarshal(apiResp.Data, &user); err != nil {
		return 0, err
	}
	return user.UserID, nil
}
