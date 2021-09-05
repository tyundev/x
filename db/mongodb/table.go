package mongodb

import (
	"fmt"

	"github.com/reiwav/x/rest"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ERR_EXIST = rest.BadRequest("exists unique in db")
)

type Table struct {
	*mgo.Collection
	Name   string
	Prefix string
	Length int
}

func NewTable(name, prefix string, length int) *Table {
	fmt.Println("DB: " + name)
	return &Table{
		Collection: NewCollection(name),
		Name:       name,
		Prefix:     prefix,
		Length:     length,
	}
}

func (t *Table) Create(model IModel) error {
	model.BeforeCreate(t.Prefix, t.Length)
	var err = t.Insert(model)
	if err != nil {
		logDB.Errorf("Create "+err.Error(), model)
	}
	return err
}

func (t *Table) UnsafeUpSertID(id string, model IModel) error {
	model.BeforeCreate(t.Prefix, t.Length)
	return t.UpSert(id, model)
}

func (t *Table) UpSert(id string, data interface{}) error {
	_, err := t.UpsertId(id, bson.M{"$set": data})
	if err != nil {
		logDB.Errorf("UpSert "+err.Error(), data)
	}
	return err
}

func (t *Table) CreateUnique(query bson.M, model IModel) error {
	count, err := t.CountWhere(query)
	if err == nil {
		if count == 0 {
			return t.Create(model)
		}
		logDB.Errorf("CreateUnique "+err.Error(), model)
		return ERR_EXIST
	}
	logDB.Errorf("CreateUnique "+err.Error(), err)
	return err
}

func (t *Table) CountWhere(query bson.M) (int, error) {
	query["deleted_at"] = 0
	var c, err = t.Find(query).Count()
	if err != nil {
		logDB.Errorf("CountWhere "+err.Error(), query)
	}
	return c, err
}

func (t *Table) FindWhere(query bson.M, result interface{}) error {
	query["deleted_at"] = 0
	var err = t.Find(query).All(result)
	if err != nil {
		logDB.Errorf("CountWhere "+err.Error(), query)
	}
	return err
}
func (t *Table) FindOne(query bson.M, result interface{}) error {
	query["deleted_at"] = 0
	var err = t.Find(query).One(result)
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}
func (t *Table) FindByID(id string, result interface{}) error {
	var err = t.FindOne(bson.M{"_id": id,
		"deleted_at": 0}, result)
	if err != nil {
		logDB.Errorf("FindByID "+err.Error(), id)
	}
	return err
}

func (t *Table) UpdateBulkModel(result map[string]IModel) error {
	var bulk = t.Bulk()
	for key, val := range result {
		val.BeforeUpdate()
		var selector = bson.M{"_id": key}
		var update = bson.M{"$set": val}
		bulk.Update(selector, update)
	}
	_, err := bulk.Run()
	if err != nil {
		logDB.Errorf("UpdateArrayBoject " + err.Error())
	}
	return err
}

func (t *Table) UpdateBulk(result map[string]interface{}) error {
	var bulk = t.Bulk()
	for key, val := range result {
		var selector = bson.M{"_id": key}
		var update = bson.M{"$set": val}
		bulk.Update(selector, update)
	}
	_, err := bulk.Run()
	if err != nil {
		logDB.Errorf("UpdateArrayBoject " + err.Error())
	}
	return err
}

func (t *Table) DeleteByID(id string) error {
	var err = t.UpdateId(id, bson.M{"$set": bson.M{"deleted_at": time.Now().Unix()}})
	if err != nil {
		logDB.Errorf("DeleteByID "+err.Error(), id)
	}
	return err
}

func (t *Table) UnsafeUpdateByID(id string, model IModel) error {
	model.BeforeUpdate()
	var err = t.UpdateId(id, bson.M{"$set": model})
	if err != nil {
		logDB.Errorf("UnsafeUpdateByID "+err.Error()+" id: "+id, model)
	}
	return err
}

func (t *Table) UpSetByID(id string, data interface{}) error {
	var err = t.UpdateId(id, bson.M{"$set": data})
	if err != nil {
		logDB.Errorf("UpSetByID "+err.Error()+" id: "+id, data)
	}
	return err
}

func (t *Table) UpdateWhere(selector interface{}, data interface{}) error {

	var err = t.Update(selector, bson.M{"$set": data})
	if err != nil {
		logDB.Errorf("UpdateWhere "+err.Error(), data)
	}
	return err
}

func (t *Table) UnsafeFindSort(queryMatch bson.M, fields []string, result interface{}) error {
	queryMatch["deleted_at"] = 0
	var err = t.Find(queryMatch).Sort(fields...).All(result)
	if err != nil {
		logDB.Errorf("UnsafeFindSort "+err.Error(), " fields: ", fields, queryMatch)
	}
	return err
}

func (t *Table) UnsafeFindSortSkip(queryMatch bson.M, fields []string, skip int, limit int, result interface{}) error {
	queryMatch["deleted_at"] = 0
	var err = t.Find(queryMatch).Sort(fields...).Skip(skip).Limit(limit).All(result)
	if err != nil {
		logDB.Errorf("UnsafeFindSort "+err.Error(), " fields: ", fields, queryMatch)
	}
	return err
}

func (t *Table) UnsafeFindSortOne(queryMatch bson.M, fields []string, result interface{}) error {
	queryMatch["deleted_at"] = 0
	var err = t.Find(queryMatch).Sort(fields...).Limit(1).One(result)
	if err != nil {
		logDB.Errorf("UnsafeFindSort "+err.Error()+" fields: ", queryMatch)
	}
	return err
}

func (t *Table) UnsafeCreateList(data []interface{}) error {
	err := t.Insert(data...)
	if err != nil {
		logDB.Errorf("UnsafeCreateAll " + err.Error())
	}
	return err
}

func (t *Table) UnsafeUpdateAll(filter bson.M, data interface{}) error {
	filter["deleted_at"] = 0
	_, err := t.UpdateAll(filter, bson.M{"$set": data})
	if err != nil {
		logDB.Errorf("UnsafeUpdateAll " + err.Error())
	}
	return err
}
