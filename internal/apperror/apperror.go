package apperror

import (
	"net/http"

	"github.com/lattiq/foundry/errors"
)

// used when the decision targets a request which is NOT in "pending" state
var ErrRequestNotPending = errors.NewAPIError("REQUEST_NOT_PENDING", http.StatusConflict)