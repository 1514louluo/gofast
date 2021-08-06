package global

import (
	"gorm.io/gorm"
	"time"
)

type MODEL struct {
	ID        uint `json:"id" form:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at" form:"created_at" gorm:"<-:create"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
