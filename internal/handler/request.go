package handler

import (
      "net/http"
      "strconv"

      "github.com/gin-gonic/gin"
      "github.com/lattiq/foundry/errors"

      "github.com/lattiq-bhuvan/access-desk/internal/dto"
      "github.com/lattiq-bhuvan/access-desk/internal/service"
)

type RequestHandler struct{ svc *service.RequestService }

func NewRequestHandler(svc *service.RequestService) *RequestHandler { return &RequestHandler{svc: svc} }

func toSvcCaller(c Caller) service.Caller { return service.Caller{UserID: c.UserID, Role: c.Role} }

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		e := errors.ErrBadRequest("invalid id")
		c.AbortWithStatusJSON(e.HttpStatus, e)
		return 0, false
	}
	return uint(id), true
}

func (h *RequestHandler) Create(c *gin.Context) {
	var body dto.CreateRequest
	if !bindJSON(c, &body) {
		return
	}
	out, err := h.svc.Create(c.Request.Context(), toSvcCaller(CallerFrom(c)), body)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *RequestHandler) List(c *gin.Context) {
	var q dto.ListRequestsQuery
	if !bindQuery(c, &q) {
		return
	}
	q.Normalize()
	rows, total, err := h.svc.List(c.Request.Context(), toSvcCaller(CallerFrom(c)), q)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"requests": rows,
		"page":     q.Page, "page_size": q.PageSize, "total": total,
	})
}

func (h *RequestHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	out, err := h.svc.Get(c.Request.Context(), toSvcCaller(CallerFrom(c)), id)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *RequestHandler) Decide(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var body dto.Decision
	if !bindJSON(c, &body) {
		return
	}
	out, err := h.svc.Decide(c.Request.Context(), toSvcCaller(CallerFrom(c)), id, body)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}