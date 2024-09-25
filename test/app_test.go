package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"E:\Golang-Shortlink\main"

	"github.com/stretchr/testify/mock"
)

const (
	expTime = 60
	longURL = "https://www.example.com"
	shortlink = "IFHzaO"
	ShortlinkInfo = `{"url":"https://www.example.com", "created_at": "2024-09-25 16:25:06.1003466 +0800 CST m=+11.493008601"}`
)

type storageMock struct{
	mock.Mock
}

var app main.App

var mockR *storageMock

func (s *storageMock) Shorten(url string, exp int64)(string, error){
	args := s.Called(url,exp)
	return args.String(0),args.Error(1)
}


func (s *storageMock) Unshorten(eid string)(string, error){
	args := s.Called(eid)
	return args.String(0),args.Error(1)
}


func (s *storageMock) ShortlinkInfo(eid string)(interface{}, error){
	args := s.Called(eid)
	return args.String(0),args.Error(1)
}


func init(){
	app = main.App{}
	mockR = new(storageMock)
	app.Initialize(&main.Env{S: mockR})
}

func TestCreateShortlink(t *testing.T){
	var jsonStr = []byte(`{
		"url": "https://www.example.com",
		"expiration_in_minutes": 60
	}`)
	req, err := http.NewRequest("POST","/api/shorten",bytes.NewBuffer(jsonStr))
	if err != nil{
		t.Fatal("Should be able to create a request.",err)
	}
	req.Header.Set("Content-Type","application/json")

	mockR.On("Shorten", longURL,int64(expTime)).Return(shortlink,nil).Once()

	rw := httptest.NewRecorder()
	app.Router.ServeHTTP(rw,req)

	if rw.Code != http.StatusCreated{
		t.Fatalf("Excepted receive %d. Got %d", http.StatusCreated,rw.Code)
	}

	resp := struct{
		Shortlink string `json: "shortlink"`
	}{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil{
		t.Fatalf("Should decode the response")
	}

	if resp.Shortlink != shortlink{
		t.Fatalf("Excepted receive %s. Got %s", shortlink,resp.Shortlink)
	}
}