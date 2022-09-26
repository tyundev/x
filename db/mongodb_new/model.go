package mongodb_new

import (
	"time"

	"github.com/tyundev/x/math"
	"github.com/tyundev/x/mlog"
)

var logDB = mlog.NewTagLog("MONGO_DB")

type IModel interface {
	BeforeCreate(prefix string)
	BeforeUpdate()
	BeforeDelete()
}
type BaseModel struct {
	ID        string `json:"id" bson:"_id"`
	CreatedAt int64  `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt int64  `json:"deleted_at" bson:"deleted_at"`
}

func (b *BaseModel) BeforeCreate(prefix string) {
	var timeNow = time.Now().Unix()
	b.ID = math.RandXID(prefix)
	if b.CreatedAt == 0 {
		b.CreatedAt = timeNow
	}
	b.UpdatedAt = timeNow
}

func (b *BaseModel) BeforeUpdate() {
	b.UpdatedAt = time.Now().Unix()
}

func (b *BaseModel) BeforeDelete() {
	b.DeletedAt = time.Now().Unix()
}
