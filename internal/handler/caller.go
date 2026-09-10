package handler

import (
      "github.com/gin-gonic/gin"
      fmw "github.com/lattiq/foundry/service/http/middleware"
      "github.com/lattiq-bhuvan/access-desk/internal/auth"
)

type Caller struct {
      UserID string
      Role   string
}

func CallerFrom(c *gin.Context) Caller {
      cl := fmw.GetClaims(c).(*auth.AccessClaims)
      return Caller{UserID: cl.UserID, Role: cl.Role}
}