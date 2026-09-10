package model

type User struct {
      ID           uint   `gorm:"primaryKey" json:"id"`
      Email        string `gorm:"uniqueIndex;not null" json:"email"`
      Name         string `gorm:"not null" json:"name"`
      Role         string `gorm:"not null" json:"role"` // requester | approver
      PasswordHash string `gorm:"not null" json:"-"`     // json:"-" — never leaks
}