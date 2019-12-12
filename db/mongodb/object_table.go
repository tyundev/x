package mongodb

import (
	"gopkg.in/mgo.v2/bson"
)

func (t *Table) CreateObj(model ModelID) error {
	model.BeforeCreate(t.Prefix, t.Length)
	var err = t.Insert(model)
	if err != nil {
		logDB.Errorf("Create "+err.Error(), model)
	}
	return err
}

func (t *Table) CreateUniqueObj(query bson.M, model ModelID) error {
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

func (t *Table) CountWhereObj(query bson.M) (int, error) {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	var c, err = t.Find(query).Count()
	if err != nil {
		logDB.Errorf("CountWhere "+err.Error(), query)
	}
	return c, err
}

func (t *Table) FindWhereObj(query bson.M, result interface{}) error {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	var err = t.Find(query).All(result)
	if err != nil {
		logDB.Errorf("CountWhere "+err.Error(), query)
	}
	return err
}
func (t *Table) FindOneObj(query bson.M, result interface{}) error {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	var err = t.Find(query).One(result)
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}
func (t *Table) FindByIDObj(id string, result interface{}) error {
	var err = t.FindId(id).One(result)
	if err != nil {
		logDB.Errorf("FindByID "+err.Error(), id)
	}
	return err
}

func (t *Table) DeleteByIDObj(id string) error {
	var err = t.UpdateId(id, bson.M{"$set": bson.M{"updated_at": 0}})
	if err != nil {
		logDB.Errorf("DeleteByID "+err.Error(), id)
	}
	return err
}

func (t *Table) UnsafeUpdateByIDObj(id string, model ModelID) error {
	model.BeforeUpdate()
	var err = t.UpdateId(id, bson.M{"$set": model})
	if err != nil {
		logDB.Errorf("UnsafeUpdateByID "+err.Error()+" id: "+id, model)
	}
	return err
}

func (t *Table) UpSetByIDObj(id string, data interface{}) error {
	var err = t.UpdateId(id, bson.M{"$set": data})
	if err != nil {
		logDB.Errorf("UpSetByID "+err.Error()+" id: "+id, data)
	}
	return err
}

func (t *Table) UpdateWhereObj(selector interface{}, data interface{}) error {
	var err = t.Update(selector, bson.M{"$set": data})
	if err != nil {
		logDB.Errorf("UpdateWhere "+err.Error(), data)
	}
	return err
}

func (t *Table) UnsafeFindSortObj(queryMatch bson.M, fields string, result interface{}) error {
	var err = t.Find(queryMatch).Sort(fields).All(result)
	if err != nil {
		logDB.Errorf("UnsafeFindSort "+err.Error()+" fields: "+fields, queryMatch)
	}
	return err
}

func (t *Table) UnsafeFindSortOneObj(queryMatch bson.M, fields []string, result interface{}) error {
	var err = t.Find(queryMatch).Sort(fields...).Limit(1).One(result)
	if err != nil {
		logDB.Errorf("UnsafeFindSort "+err.Error()+" fields: ", queryMatch)
	}
	return err
}
