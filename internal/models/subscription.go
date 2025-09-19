package models

import (
    "time"
    "github.com/google/uuid"
)

type Subscription struct {
    ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    ServiceName string     `gorm:"not null"`
    Price       int        `gorm:"not null"`
    UserID      uuid.UUID  `gorm:"type:uuid;not null;index"`
    StartDate   time.Time  `gorm:"not null"`
    EndDate     *time.Time
}
