package service

import (
	"net/http"
	"weixin_backend/utils"

	"github.com/gin-gonic/gin"
)

const (
	None = 0
	Bind = 1 << iota
	BindUri
	BindQuery
)

type Service interface {
	Handle(c *gin.Context) (any, error)
}

func HandlerNoBind(s Service) gin.HandlerFunc {
	return HandlerWithBindType(s, None)
}

func HandlerBind(s Service) gin.HandlerFunc {
	return HandlerWithBindType(s, Bind)
}

func HandlerBindUri(s Service) gin.HandlerFunc {
	return HandlerWithBindType(s, BindUri)
}

func HandlerBindQuery(s Service) gin.HandlerFunc {
	return HandlerWithBindType(s, BindQuery)
}

func HandlerWithBindType(s Service, bindType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		if bindType&BindUri != 0 {
			if err = c.ShouldBindUri(s); err != nil {
				c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
				return
			}
		}
		if bindType&Bind != 0 {
			if err = c.ShouldBind(s); err != nil {
				c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
				return
			}
		}
		if bindType&BindQuery != 0 {
			if err = c.ShouldBindQuery(s); err != nil {
				c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
				return
			}
		}

		res, err := s.Handle(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		} else {
			c.JSON(http.StatusOK, utils.Response(res))
		}
	}
}
