package main_test

import (
	"Golang-Shortlink/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

const (
	expTime   = 60
	longURL   = "https://www.example.com"
	shortlink = "IFHzaO"
)

type storageMock struct {
	mock.Mock
}

type URLDetail struct {
	URL                 string        `json:"url"`
	CreateAt            string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

var a app.App

var mockR *storageMock

func (s *storageMock) Shorten(url string, exp int64) (string, error) {
	args := s.Called(url, exp)
	return args.String(0), args.Error(1)
}

func (s *storageMock) Unshorten(eid string) (string, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func (s *storageMock) ShortlinkInfo(eid string) (interface{}, error) {
	args := s.Called(eid)
	return args.Get(0), args.Error(1)
}

func init() {
	a = app.App{}
	mockR = new(storageMock)
	a.Initialize(&app.Env{S: mockR})
}

func TestCreateShortlink(t *testing.T) {
	var jsonStr = []byte(`{
		"url": "https://www.example.com",
		"expiration_in_minutes": 60
	}`)
	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal("Should be able to create a request.", err)
	}
	req.Header.Set("Content-Type", "application/json")

	mockR.On("Shorten", longURL, int64(expTime)).Return(shortlink, nil).Once()

	rw := httptest.NewRecorder()
	a.Router.ServeHTTP(rw, req)

	if rw.Code != http.StatusAccepted {
		t.Fatalf("Excepted receive %d. Got %d", http.StatusAccepted, rw.Code)
	}

	resp := struct {
		Shortlink string `json:"shortlink"`
	}{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("Should decode the response")
	}

	if resp.Shortlink != shortlink {
		t.Fatalf("Excepted receive %s. Got %s", shortlink, resp.Shortlink)
	}
}

func TestRedirect(t *testing.T) {
	r := fmt.Sprintf("/%s", shortlink)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Fatal("Should be able to create a request.", err)
	}

	mockR.On("Unshorten", shortlink).Return(longURL, nil).Once()

	rw := httptest.NewRecorder()
	a.Router.ServeHTTP(rw, req)

	if rw.Code != http.StatusTemporaryRedirect {
		t.Fatalf("Excepted receive %d. Got %d", http.StatusTemporaryRedirect, rw.Code)
	}
}

func TestGetShortlinkInfo(t *testing.T) {
	// 构建请求，带上 shortlink 参数
	req, err := http.NewRequest("GET", "/api/info?shortlink="+shortlink, nil)
	if err != nil {
		t.Fatal("Should be able to create a request.", err)
	}

	// 设置期望的返回值，使用 json.Marshal 生成合法的 JSON 字符串
	shortlinkInfo := &URLDetail{
		URL:                 longURL,
		CreateAt:            time.Now().String(),
		ExpirationInMinutes: expTime * time.Second,
	}

	// shortlinkInfoJSON, err := json.Marshal(shortlinkInfo)
	// if err != nil {
	// 	t.Fatalf("Failed to marshal shortlinkInfo: %v", err)
	// }

	// 设置 mock 行为，返回 JSON 字符串
	mockR.On("ShortlinkInfo", shortlink).Return(shortlinkInfo, nil).Once()

	// 使用 httptest 包进行测试
	rw := httptest.NewRecorder()
	a.Router.ServeHTTP(rw, req)

	// 打印响应体
	fmt.Println("rw.Body: ", rw.Body)

	// 检查返回的状态码是否为 200 OK
	if rw.Code != http.StatusOK {
		t.Fatalf("Expected receive %d. Got %d", http.StatusOK, rw.Code)
	}

	// 检查响应体是否包含期望的数据
	var respBody URLDetail
	if err := json.NewDecoder(rw.Body).Decode(&respBody); err != nil {
		t.Fatalf("Should decode the response: %v", err)
	}

	fmt.Println("resp type:", reflect.TypeOf(respBody))

	// 逐字段比较，而不是直接比较结构体
	if respBody.URL != longURL {
		t.Fatalf("Expected receive %v. Got %v", longURL, respBody.URL)
	}
}
