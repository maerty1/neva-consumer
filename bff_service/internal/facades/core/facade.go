package core

import (
	"github.com/gin-gonic/gin"
)

type Facade interface {
	GetCurrentWeather(c *gin.Context)
}

var _ Facade = (*facade)(nil)

type facade struct {
}

func NewFacade() Facade {
	return &facade{}
}
