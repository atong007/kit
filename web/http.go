package web

import (
	"github.com/gin-gonic/gin"
)

func NewServer(addr string) *gin.Engine {
	e := gin.Default()
	return e
}
