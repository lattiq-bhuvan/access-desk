package handler

import (
      "github.com/gin-gonic/gin"
      "github.com/lattiq/foundry/errors"
)

type Caller struct {
	UserID string
	Role   string // "requester" | "approver"
}

const callerKey = "caller"

// TempAuth reads X-User-Id / X-Role headers. M4 replaces this with jwt middleware
func TempAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract uid and role from http request headers
		uid := c.GetHeader("X-User-Id")
		role := c.GetHeader("X-Role")
		if uid == "" || (role != "requester" && role != "approver") {
			e := errors.ErrUnauthorized("missing or invalid X-User-Id / X-Role")
			c.AbortWithStatusJSON(e.HttpStatus, e)
			return
		}

		// store the validated(not authenticated) headers in gin context to pass it downstream
		c.Set(callerKey, Caller{UserID: uid, Role: role})
		c.Next()
	}
}

func CallerFrom(c *gin.Context) Caller {
    return c.MustGet(callerKey).(Caller)
}