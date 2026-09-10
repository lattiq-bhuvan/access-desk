package service

import (
	stderrors "errors"

	"context"

	"gorm.io/gorm"

	fErr "github.com/lattiq/foundry/errors"

	"github.com/lattiq-bhuvan/access-desk/internal/apperror"
	"github.com/lattiq-bhuvan/access-desk/internal/dto"
	"github.com/lattiq-bhuvan/access-desk/internal/model"
	"github.com/lattiq-bhuvan/access-desk/internal/store"
)

type Caller struct {
	UserID string
	Role   string
}

type RequestService struct{ st *store.Store }

func NewRequestService(st *store.Store) *RequestService { return &RequestService{st: st} }

func (s *RequestService) Create(ctx context.Context, c Caller, in dto.CreateRequest) (*model.AccessRequest, error) {
	r := &model.AccessRequest{
		DatasetID: in.DatasetID,
		Requester: c.UserID,
		Reason:    in.Reason,
		Status:    "PENDING",
	}
	err := s.st.CreateRequest(ctx, r)
	switch {
	case stderrors.Is(err, gorm.ErrDuplicatedKey):
		return nil, fErr.ErrResourceAlreadyExists("you already requested this dataset")
	case stderrors.Is(err, gorm.ErrForeignKeyViolated):
		return nil, fErr.ErrBadRequest("dataset does not exist")
	case err != nil:
		return nil, err
	}
    return r, nil
}

func (s *RequestService) Get(ctx context.Context, c Caller, id uint) (*model.AccessRequest, error) {
	r, err := s.st.GetRequest(ctx, id)
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fErr.ErrResourceNotFound("request not found")
	}
	if err != nil {
		return nil, err
	}
	if c.Role != "approver" && r.Requester != c.UserID {
		return nil, fErr.ErrAccessDenied("not your request")
	}
	return r, nil
}

func (s *RequestService) List(ctx context.Context, c Caller, q dto.ListRequestsQuery) ([]model.AccessRequest, int64, error) {
      onlyRequester := ""
      if c.Role != "approver" {
              onlyRequester = c.UserID // requester sees only their own
      }
      return s.st.ListRequests(ctx, onlyRequester, q.Status, q.Page, q.PageSize)
}

func (s *RequestService) Decide(ctx context.Context, c Caller, id uint, in dto.Decision) (*model.AccessRequest, error) {
      if c.Role != "approver" {
              return nil, fErr.ErrAccessDenied("only approvers can decide requests")
      }

      newStatus := map[string]string{"APPROVE": "APPROVED", "REJECT": "REJECTED"}[in.Decision]

      check := func(cur *model.AccessRequest) error {
              if cur.Requester == c.UserID {
                      return fErr.ErrAccessDenied("cannot decide your own request")
              }
              if cur.Status != "PENDING" {
                      // §2.2: terminal request -> 409 RESOURCE_ALREADY_EXISTS
                      return fErr.ErrResourceAlreadyExists("request has already been decided")
              }
              return nil
      }

      dec := &model.Decision{Decider: c.UserID, Decision: in.Decision, Note: in.Note}
      err := s.st.DecideRequest(ctx, id, check, newStatus, dec)
      if stderrors.Is(err, gorm.ErrRecordNotFound) {
              return nil, fErr.ErrResourceNotFound("request not found")
      }
      if err != nil {
              return nil, err // an *APIError from check() passes through GetServiceError unchanged
      }
      return s.st.GetRequest(ctx, id)
}

// keep apperr import alive for when you use ErrRequestNotPending
var _ = apperror.ErrRequestNotPending