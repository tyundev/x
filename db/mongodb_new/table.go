package mongodb_new

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Table struct {
	*mongo.Collection
	Prefix string
}

func NewTable(name, prefix string, db *mongo.Database) *Table {
	fmt.Println("DB", name)
	return &Table{
		Collection: db.Collection(name),
		Prefix:     prefix,
	}
}

func (t *Table) Create(model IModel) error {
	ctx := context.Background()
	model.BeforeCreate(t.Prefix)
	var _, err = t.InsertOne(ctx, model)
	if err != nil {
		logDB.Errorf("Create table "+t.Name()+": "+err.Error(), model)
	}
	return err
}

func (t *Table) Update(id string, model IModel) error {
	ctx := context.Background()
	model.BeforeUpdate()
	var _, err = t.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": 0}, bson.M{"$set": model})
	if err != nil {
		logDB.Errorf("Update table "+t.Name()+": "+err.Error(), model)
	}
	return err
}

func (t *Table) Delete(id string, model IModel) error {
	ctx := context.Background()
	model.BeforeDelete()
	var _, err = t.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": 0}, bson.M{"$set": model})
	if err != nil {
		logDB.Errorf("Delete table "+t.Name()+": "+err.Error(), model)
	}
	return err
}

func (t *Table) UnsafeUpdateByID(id string, v interface{}) error {
	ctx := context.Background()
	var _, err = t.UpdateOne(ctx,
		bson.M{"deleted_at": 0,
			"_id": id},
		bson.M{"$set": v})
	if err != nil {
		logDB.Errorf("UnsafeUpdateByID table "+t.Name()+": "+err.Error(), v)
	}
	return err
}

func (t *Table) CreateMany(v []interface{}) ([]interface{}, error) {
	ctx := context.Background()
	var res, err = t.InsertMany(ctx, v)
	var ids []interface{}
	if err != nil {
		logDB.Errorf("UnsafeUpdateByID table "+t.Name()+": "+err.Error(), v)
	}
	if res != nil {
		ids = res.InsertedIDs
	}
	return ids, err
}

func (t *Table) SelectOne(filter bson.M, v interface{}) error {
	ctx := context.Background()
	filter["deleted_at"] = 0
	var err = t.FindOne(ctx, filter).Decode(v)
	return err
}

func (t *Table) SelectOneWithFields(filter bson.M, v interface{}, fields bson.M) error {
	ctx := context.Background()
	filter["deleted_at"] = 0

	var opts = options.FindOne().SetProjection(fields)
	var err = t.FindOne(ctx, filter, opts).Decode(v)
	return err
}

func (t *Table) SelectManyWithFields(filter bson.M, v interface{}, fields bson.M) error {
	ctx := context.Background()
	filter["deleted_at"] = 0

	var opts = options.Find().SetProjection(fields)
	var cur, err = t.Find(ctx, filter, opts)
	if err != nil {
		cur.Close(ctx)
		return err
	}
	err = cur.All(ctx, v)
	return err
}

func (t *Table) SelectByID(id string, v interface{}) error {
	ctx := context.Background()
	var filter = bson.M{
		"deleted_at": 0,
		"_id":        id,
	}
	var err = t.FindOne(ctx, filter).Decode(v)

	return err
}

func (t *Table) SelectMany(filter bson.M, v interface{}) error {
	ctx := context.Background()
	filter["deleted_at"] = 0
	var cur, err = t.Find(ctx, filter)
	if err != nil {
		cur.Close(ctx)
		return err
	}
	err = cur.All(ctx, v)
	return err
}

func (t *Table) UpdateAll(filter bson.M, update interface{}) error {
	ctx := context.Background()
	filter["deleted_at"] = 0
	var _, err = t.UpdateMany(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return err
}

func (t *Table) SelectAndSort(filter bson.M, sort interface{}, skip, limit int64, res interface{}) error {
	ctx := context.Background()
	filter["deleted_at"] = 0
	var opts = options.Find()
	if sort != nil {
		opts.SetSort(sort)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	var cur, err = t.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	err = cur.All(ctx, res)
	return err
}

func (t *Table) Pipe(pipeline mongo.Pipeline, res interface{}) error {
	ctx := context.Background()
	var cur, err = t.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	err = cur.All(ctx, res)
	return err
}

func (t *Table) Count(filter bson.M) (int64, error) {
	ctx := context.Background()
	return t.CountDocuments(ctx, filter)
}
