package mongodb

import (
	"time"
	"x/math"
)

type IModel interface {
	BeforeCreate(prefix string, length int)
	BeforeUpdate()
	BeforeDelete()
}
type BaseModel struct {
	ID        string `json:"id" bson:"_id"`
	CreatedAt int64  `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at,omitempty"`
	DeletedAt int64  `json:"deleted_at" bson:"deleted_at"`
}

func (b *BaseModel) BeforeCreate(prefix string, length int) {
	b.ID = math.RandString(prefix, length)
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
}

func (b *BaseModel) BeforeUpdate() {
	b.UpdatedAt = time.Now().Unix()
}

func (b *BaseModel) BeforeDelete() {
	b.DeletedAt = time.Now().Unix()
}
