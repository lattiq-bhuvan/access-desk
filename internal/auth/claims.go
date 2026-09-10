package auth

import "errors"

type AccessClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"` // "requester" | "approver"
}

func (c *AccessClaims) ValidateClaims() error {
	if c.UserID == "" {
		return errors.New("missing user_id")
	}
	if c.Role != "requester" && c.Role != "approver" {
		return errors.New("invalid role")
	}
	return nil
}

func ClaimsFactory() any {
	return &AccessClaims{}
}