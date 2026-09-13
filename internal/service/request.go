package service

import (
	stderrors "errors"
	"strconv"

	"context"

	"gorm.io/gorm"

	fErr "github.com/lattiq/foundry/errors"
	"github.com/lattiq/foundry/o11y"
	"github.com/lattiq/foundry/o11y/tracing"

	"github.com/lattiq-bhuvan/access-desk/internal/apperror"
	"github.com/lattiq-bhuvan/access-desk/internal/dto"
	"github.com/lattiq-bhuvan/access-desk/internal/model"
	"github.com/lattiq-bhuvan/access-desk/internal/store"
)

type Caller struct {
	UserID string
	Role   string
}

// id parses the caller's JWT-carried user id (kept as a string in the token
// to dodge JSON-number precision issues) back into the uint used as the
// real DB foreign key. It only fails for a hand-crafted/forged token, so a
// parse failure is treated as an auth error, not a 500.
func (c Caller) id() (uint, error) {
	id, err := strconv.ParseUint(c.UserID, 10, 64)
	if err != nil {
		return 0, fErr.ErrUnauthorized("invalid caller id")
	}
	return uint(id), nil
}

type RequestService struct {
	st       *store.Store
	observer o11y.Observer
}

func NewRequestService(st *store.Store) *RequestService {
	return &RequestService{
		st:       st,
		observer: o11y.NewObserver("request-service"),
	}
}

func (s *RequestService) Create(ctx context.Context, c Caller, in dto.CreateRequest) (*model.AccessRequest, error) {
	ctx, span := s.observer.StartSpan(ctx)
	defer span.End()

	requesterID, err := c.id()
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	r := &model.AccessRequest{
		DatasetID:   in.DatasetID,
		RequesterID: requesterID,
		Reason:      in.Reason,
		Status:      "PENDING",
	}
	err = s.st.CreateRequest(ctx, r)
	switch {
	case stderrors.Is(err, gorm.ErrDuplicatedKey):
		{
			tracing.RecordError(span, err)
			return nil, fErr.ErrResourceAlreadyExists("you already requested this dataset")
		}
	case stderrors.Is(err, gorm.ErrForeignKeyViolated):
		{
			tracing.RecordError(span, err)
			return nil, fErr.ErrBadRequest("dataset does not exist")
		}
	case err != nil:
		{
			tracing.RecordError(span, err)
			return nil, err
		}
	}
	// Create doesn't preload the Requester association (nothing to preload —
	// we just wrote the row); re-fetch through GetRequest so the response
	// still carries a human-readable requester email like every other route.
	return s.st.GetRequest(ctx, r.ID)
}

func (s *RequestService) Get(ctx context.Context, c Caller, id uint) (*model.AccessRequest, error) {
	ctx, span := s.observer.StartSpan(ctx)
	defer span.End()

	callerID, err := c.id()
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	r, err := s.st.GetRequest(ctx, id)
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		tracing.RecordError(span, err)
		return nil, fErr.ErrResourceNotFound("request not found")
	}
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}
	if c.Role != "approver" && r.RequesterID != callerID {
		err := fErr.ErrAccessDenied("not your request")
		tracing.RecordError(span, err)
		return nil, err
	}
	return r, nil
}

func (s *RequestService) List(ctx context.Context, c Caller, q dto.ListRequestsQuery) ([]model.AccessRequest, int64, error) {
	ctx, span := s.observer.StartSpan(ctx)
	defer span.End()

	var onlyRequester uint
	if c.Role != "approver" {
		id, err := c.id()
		if err != nil {
			tracing.RecordError(span, err)
			return nil, 0, err
		}
		onlyRequester = id // requester sees only their own
	}
	return s.st.ListRequests(ctx, onlyRequester, q.Status, q.Page, q.PageSize)
}

func (s *RequestService) Decide(ctx context.Context, c Caller, id uint, in dto.Decision) (*model.AccessRequest, error) {
	ctx, span := s.observer.StartSpan(ctx)
	defer span.End()

	tracing.AddSpanAttributes(span, map[string]any{
		"request.id":     id,
		"caller.role":    c.Role,
		"decision.value": in.Decision,
	})

	if c.Role != "approver" {
		return nil, fErr.ErrAccessDenied("only approvers can decide requests")
	}

	callerID, err := c.id()
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	newStatus := map[string]string{"APPROVE": "APPROVED", "REJECT": "REJECTED"}[in.Decision]

	check := func(cur *model.AccessRequest) error {
		if cur.RequesterID == callerID {
			return fErr.ErrAccessDenied("cannot decide your own request")
		}
		if cur.Status != "PENDING" {
			// §2.2: terminal request -> 409 RESOURCE_ALREADY_EXISTS
			return fErr.ErrResourceAlreadyExists("request has already been decided")
		}
		return nil
	}

	dec := &model.Decision{DeciderID: callerID, Decision: in.Decision, Note: in.Note}
	err = s.st.DecideRequest(ctx, id, check, newStatus, dec)
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		tracing.RecordError(span, err)
		return nil, fErr.ErrResourceNotFound("request not found")
	}
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err // an *APIError from check() passes through GetServiceError unchanged
	}
	return s.st.GetRequest(ctx, id)
}

// keep apperr import alive for when you use ErrRequestNotPending
var _ = apperror.ErrRequestNotPending
