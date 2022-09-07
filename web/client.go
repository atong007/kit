package web

import (
	"net/http"
	"sync"
	"time"
)

var (
	client *http.Client
	once   sync.Once
)

func Instance() *http.Client {
	once.Do(func() {
		client = &http.Client{
			Timeout: time.Second * 60,
		}
	})
	return client
}
