package web

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type Data struct {
	Val int `json:"val"`
}

func TestRequest_Get(t *testing.T) {

	type args struct {
		respBytes []byte
	}

	tests := []struct {
		name     string
		args     args
		wantCode int
		wantResp any
		wantErr  bool
	}{
		{
			"success",
			args{
				respBytes: []byte(`{"code": 200, "msg": "ok", "data": {"val": 1}}`),
			},
			200,
			Response[Data]{
				Code: 200,
				Msg:  "ok",
				Data: Data{Val: 1},
			},
			false,
		},
		{
			"success with code 400",
			args{
				respBytes: []byte(`{"code": 400, "msg": "request failed", "data": {}}`),
			},
			400,
			Response[Data]{
				Code: 400,
				Msg:  "request failed",
				Data: Data{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.wantCode)
				w.Write(tt.args.respBytes)
			}))
			defer ts.Close()

			gotCode, gotResp, err := NewReq[Response[Data]]().Get(ts.URL)
			assert.Equal(t, tt.wantCode, gotCode)
			assert.Equal(t, tt.wantResp, gotResp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequest_Post(t *testing.T) {

	type args struct {
		body      interface{}
		respBytes []byte
	}

	tests := []struct {
		name     string
		args     args
		wantCode int
		wantResp any
		wantErr  bool
	}{
		{
			"success",
			args{
				body: map[string]interface{}{
					"key1": "value2",
					"key2": 2,
				},
				respBytes: []byte(`{"code": 200, "msg": "ok", "data": {"val": 1}}`),
			},
			200,
			Response[Data]{
				Code: 200,
				Msg:  "ok",
				Data: Data{Val: 1},
			},
			false,
		},
		{
			"success with code 400",
			args{
				body: map[string]interface{}{
					"key1": "value2",
					"key2": 2,
				},
				respBytes: []byte(`{"code": 400, "msg": "request failed", "data": {}}`),
			},
			400,
			Response[Data]{
				Code: 400,
				Msg:  "request failed",
				Data: Data{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bs, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err)
				bodyBytes, err := json.Marshal(tt.args.body)
				require.NoError(t, err)
				require.Equal(t, bodyBytes, bs)
				w.WriteHeader(tt.wantCode)
				w.Write(tt.args.respBytes)
			}))
			defer ts.Close()

			gotCode, gotResp, err := NewReq[Response[Data]]().SetBody(tt.args.body).Post(ts.URL)
			assert.Equal(t, tt.wantCode, gotCode)
			assert.Equal(t, tt.wantResp, gotResp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequest_SetHeader(t *testing.T) {

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"success",
			args{
				key:   "content-type",
				value: "application/json",
			},
			"application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReq[string]()
			r = r.SetHeader(tt.args.key, tt.args.value)
			//query := r.req.URL.Query().Encode()
			assert.Equal(t, tt.want, r.req.Header.Get(tt.args.key))
		})
	}
}

func TestRequest_SetQueryParams(t *testing.T) {

	type args struct {
		params map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"success",
			args{
				params: map[string]string{
					"username": "charlie",
				},
			},
			map[string]string{
				"username": "charlie",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReq[string]()
			r = r.SetQueryParams(tt.args.params)
			assert.Equal(t, tt.want, r.queryParams)
		})
	}
}

func TestRequest_SetBody(t *testing.T) {

	type args struct {
		body interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"success",
			args{
				body: map[string]string{
					"username": "charlie",
				},
			},
			map[string]string{
				"username": "charlie",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReq[string]()
			r = r.SetBody(tt.args.body)
			assert.Equal(t, tt.want, r.reqBody)
		})
	}
}
