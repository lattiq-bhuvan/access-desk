package model

import "time"

type Dataset struct {
      ID          uint      `gorm:"primaryKey" json:"id"`
      Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
      Name        string    `gorm:"not null" json:"name"`
      Description string    `json:"description"`
      CreatedAt   time.Time `json:"created_at"`
}

type AccessRequest struct {
      ID        uint      `gorm:"primaryKey" json:"id"`
      DatasetID uint      `gorm:"not null;uniqueIndex:idx_ds_requester" json:"dataset_id"`
      Dataset   Dataset   `gorm:"foreignKey:DatasetID" json:"dataset"`
      Requester string    `gorm:"not null;uniqueIndex:idx_ds_requester" json:"requester"`      // user id/email
      Reason    string    `json:"reason"`
      Status    string    `gorm:"not null;default:PENDING" json:"status"` // PENDING|APPROVED|REJECTED
      CreatedAt time.Time `json:"created_at"`
      UpdatedAt time.Time `json:"updated_at"`
}

type Decision struct {
      ID              uint      `gorm:"primaryKey" json:"id"`
      AccessRequestID uint      `gorm:"not null;uniqueIndex" json:"access_request_id"` // one decision per request
      Decider         string    `gorm:"not null" json:"decider"`
      Decision        string    `gorm:"not null" json:"decision"` // APPROVE|REJECT
      Note            string    `json:"note"`
      CreatedAt       time.Time `json:"created_at"`
}