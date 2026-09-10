package dto

type CreateRequest struct {
	DatasetID uint `json:"dataset_id" validate:"required"`
	Reason string `json:"reason" validate:"required,min=1,max=500"`
}

type Decision struct {
	Decision string `json:"decision" validate:"required,oneof=APPROVE REJECT"`
	Note string `json:"note" validate:"max=500"`
}

// ListRequestsQuery binds from the query string: ?status=PENDING&page=1&page_size=20
type ListRequestsQuery struct {
      Status   string `form:"status" validate:"omitempty,oneof=PENDING APPROVED REJECTED"`
      Page     int    `form:"page" validate:"omitempty,min=1"`
      PageSize int    `form:"page_size" validate:"omitempty,min=1,max=100"`
}

// provide default values for pagination if omitted in query string
func (q *ListRequestsQuery) Normalize() {
      if q.Page == 0 {
              q.Page = 1
      }
      if q.PageSize == 0 {
              q.PageSize = 20
      }
}