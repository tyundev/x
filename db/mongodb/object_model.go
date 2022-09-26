package mongodb

import (
	"time"

	"github.com/tyundev/x/math"

	"gopkg.in/mgo.v2/bson"
)

type ModelID interface {
	BeforeCreateObj(prefix string, length int)
	BeforeUpdateObj()
	BeforeDeleteObj()
	//GetID()
}

type BaseModelObj struct {
	OID       bson.ObjectId `json:"-" bson:"_id"`
	ID        string        `json:"id" bson:"id_str,omitempty"`
	CreatedAt int64         `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt int64         `json:"updated_at" bson:"updated_at,omitempty"`
}

func (b *BaseModelObj) BeforeCreateObj(prefix string, length int) {
	b.ID = math.RandString(prefix, length)
	b.OID = bson.ObjectId(b.ID)
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
}

func (b *BaseModelObj) BeforeUpdateObj() {
	b.UpdatedAt = time.Now().Unix()
}

func (b *BaseModelObj) BeforeDeleteObj() {
	b.UpdatedAt = 0
}

func (b *BaseModelObj) GetID() string {
	b.ID = b.OID.Hex()
	return b.ID
}
