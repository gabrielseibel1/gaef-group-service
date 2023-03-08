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

func (h Handler) AuthMiddleware() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) OnlyLeadersMiddleware() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) CreateGroupHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) ReadAllGroupsHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) ReadGroupHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) UpdateGroupHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) DeleteGroupHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) ReadMembersHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) AddMemberHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) DeleteMemberHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) ReadLeadersHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) AddLeadersHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

func (h Handler) DeleteLeaderHandler() gin.HandlerFunc {
	panic("not implemented") // TODO: implement me
}

type Service interface {
	Serve()
}
