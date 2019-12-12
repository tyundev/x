package mongodb

import (
	"time"
	"x/math"

	"gopkg.in/mgo.v2/bson"
)

type ModelID interface {
	BeforeCreate(prefix string, length int)
	BeforeUpdate()
	BeforeDelete()
}

type BaseModelObj struct {
	OID       bson.ObjectId `json:"-" bson:"_id"`
	ID        string        `json:"id" bson:"id_str,omitempty"`
	CreatedAt int64         `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt int64         `json:"updated_at" bson:"updated_at,omitempty"`
}

func (b *BaseModelObj) BeforeCreate(prefix string, length int) {
	b.ID = math.RandString(prefix, length)
	b.OID = bson.ObjectId(b.ID)
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
}

func (b *BaseModelObj) BeforeUpdate() {
	b.UpdatedAt = time.Now().Unix()
}

func (b *BaseModelObj) BeforeDelete() {
	b.UpdatedAt = 0
}
