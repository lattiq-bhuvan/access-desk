package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lattiq-bhuvan/access-desk/internal/auth"
	"github.com/lattiq-bhuvan/access-desk/internal/dto"
)

type AuthHandler struct{ svc *auth.Service }

func NewAuthHandler(svc *auth.Service) (h *AuthHandler){
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body dto.LoginBody
	if !bindJSON(c, &body) { return }
	res, err := h.svc.Login(c.Request.Context(), body)
	if err != nil { fail(c, err); return }
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var body dto.RefreshBody
	if !bindJSON(c, &body) { return }
	res, err := h.svc.Refresh(c.Request.Context(), body)
	if err != nil { fail(c, err); return }
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(c *gin.Context) { c.Status(http.StatusNoContent) } // stub per PDF

func (h *AuthHandler) Me(c *gin.Context) {
	caller := CallerFrom(c)
	u, err := h.svc.Me(c.Request.Context(), caller.UserID)
	if err != nil { fail(c, err); return }
	c.JSON(http.StatusOK, dto.MeResponse{
		ID:    strconv.Itoa(int(u.ID)),
		Email: u.Email,
		Name:  u.Name,
		Roles: []string{u.Role},
	})
}

// PATCH /v1/users/me/settings — webtk calls it; stub it.
func (h *AuthHandler) UpdateSettings(c *gin.Context) { c.Status(http.StatusNoContent) }