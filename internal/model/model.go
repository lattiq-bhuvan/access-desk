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
      ID          uint      `gorm:"primaryKey" json:"id"`
      DatasetID   uint      `gorm:"not null;uniqueIndex:idx_ds_requester" json:"dataset_id"`
      Dataset     Dataset   `gorm:"foreignKey:DatasetID" json:"dataset"`
      RequesterID uint      `gorm:"not null;uniqueIndex:idx_ds_requester" json:"requester_id"`
      Requester   User      `gorm:"foreignKey:RequesterID" json:"-"` // preloaded for display only
      // RequesterEmail is not a DB column (gorm:"-"): the store populates it
      // from the preloaded Requester association after every query. The JSON
      // key stays "requester" so the frontend, which already renders this
      // field as display text, needed no changes.
      RequesterEmail string    `gorm:"-" json:"requester"`
      Reason         string    `json:"reason"`
      Status         string    `gorm:"not null;default:PENDING" json:"status"` // PENDING|APPROVED|REJECTED
      CreatedAt      time.Time `json:"created_at"`
      UpdatedAt      time.Time `json:"updated_at"`
}

type Decision struct {
      ID              uint      `gorm:"primaryKey" json:"id"`
      AccessRequestID uint      `gorm:"not null;uniqueIndex" json:"access_request_id"` // one decision per request
      DeciderID       uint      `gorm:"not null" json:"decider_id"`
      Decider         User      `gorm:"foreignKey:DeciderID" json:"-"` // preloaded for display only
      // DeciderEmail is not a DB column (gorm:"-"): populated from the
      // preloaded Decider association whenever a Decision is read back.
      // The JSON key stays "decider" — same pattern as AccessRequest.RequesterEmail.
      DeciderEmail string    `gorm:"-" json:"decider"`
      Decision     string    `gorm:"not null" json:"decision"` // APPROVE|REJECT
      Note         string    `json:"note"`
      CreatedAt    time.Time `json:"created_at"`
}