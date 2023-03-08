package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func New(service Service) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.service.Serve()
	}
}

type Service interface {
	Serve()
}
