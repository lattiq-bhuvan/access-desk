package handler

import (
      "github.com/gin-gonic/gin"
      "github.com/lattiq/foundry/errors"
)

// bindJSON binds + validates the body.
// handles json binding errors
func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		apiErr := errors.GetBindingError(c.Request.Context(), err)
		c.AbortWithStatusJSON(apiErr.HttpStatus, apiErr)
		return false
	}
	return true
}

// bindQuery binds + validates the query parameters
// handles query binding errors
func bindQuery(c *gin.Context, dst any) bool {
	if err := c.ShouldBindQuery(dst); err != nil {
		apiErr := errors.GetBindingError(c.Request.Context(), err)
		c.AbortWithStatusJSON(apiErr.HttpStatus, apiErr)
		return false
	}
	return true
}

// fail converts any service-layer error into a structured response.
func fail(c *gin.Context, err error) {
      apiErr := errors.GetServiceError(c.Request.Context(), err)
      c.AbortWithStatusJSON(apiErr.HttpStatus, apiErr)
}