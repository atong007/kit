package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var defaultClient = &http.Client{
	Timeout: time.Second * 60,
}

type Request[T any] struct {
	req         *http.Request
	queryParams map[string]string
	reqBody     interface{}
	resp        T
}

func NewReq[T any]() *Request[T] {
	req := &http.Request{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}
	req.Header.Set("Content-Type", "application/json")
	return &Request[T]{req: req}
}

func (r *Request[T]) SetHeader(key, value string) *Request[T] {
	r.req.Header.Set(key, value)
	return r
}

func (r *Request[T]) SetQueryParams(p map[string]string) *Request[T] {
	r.queryParams = p
	return r
}

func (r *Request[T]) SetBody(b interface{}) *Request[T] {
	r.reqBody = b
	return r
}

func (r *Request[T]) Get(urlStr string) (code int, resp T, err error) {
	return r.exec(http.MethodGet, urlStr)
}

func (r *Request[T]) Post(urlStr string) (code int, resp T, err error) {
	return r.exec(http.MethodPost, urlStr)
}

func (r *Request[T]) Patch(urlStr string) (code int, resp T, err error) {
	return r.exec(http.MethodPatch, urlStr)
}

func (r *Request[T]) Put(urlStr string) (code int, resp T, err error) {
	return r.exec(http.MethodPut, urlStr)
}

func (r *Request[T]) Delete(urlStr string) (code int, resp T, err error) {
	return r.exec(http.MethodDelete, urlStr)
}

func (r *Request[T]) exec(method, urlStr string) (code int, resp T, err error) {
	r.req.Method = method
	switch method {
	case http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete:
		if err = r.parseBody(); err != nil {
			return
		}
	}
	if r.queryParams != nil {
		var builder strings.Builder
		builder.WriteString("?")
		for k, v := range r.queryParams {
			builder.WriteString(fmt.Sprintf("%s=%s&", k, v))
		}
		urlStr += strings.TrimRight(builder.String(), "&")
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	r.req.URL = u
	r.req.Host = u.Host

	rb, err := defaultClient.Do(r.req)
	if err != nil {
		return
	}
	defer rb.Body.Close()

	code = rb.StatusCode

	err = json.NewDecoder(rb.Body).Decode(&resp)
	if err != nil {
		return
	}
	return
}

func (r *Request[T]) parseBody() error {
	if r.reqBody != nil {
		var rc io.ReadCloser
		b, err := json.Marshal(r.reqBody)
		if err != nil {
			return err
		}

		var body io.Reader = bytes.NewBuffer(b)
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
			rc = io.NopCloser(body)
		}
		r.req.Body = rc

		v := body.(*bytes.Buffer)
		r.req.ContentLength = int64(v.Len())
		buf := v.Bytes()
		r.req.GetBody = func() (io.ReadCloser, error) {
			br := bytes.NewReader(buf)
			return io.NopCloser(br), nil
		}

	}

	if r.req.GetBody != nil && r.req.ContentLength == 0 {
		r.req.Body = http.NoBody
		r.req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
	}
	return nil
}
