package scheduler

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoJobDB struct {
	conn     *mongo.Client
	dbName   string
	collName string
}

func newMongo(
	conn *mongo.Client,
	dbName,
	collName string) JobDatabase {
	return &mongoJobDB{
		conn,
		dbName,
		collName,
	}
}

func (db *mongoJobDB) InitJobDB() (err error) {
	expiredIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "next_run_at", Value: 1},
			{Key: "status", Value: 1},
		},
	}

	statusIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
		},
	}

	nameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
	}

	_, err = db.conn.
		Database(db.dbName).
		Collection(db.collName).
		Indexes().
		CreateMany(context.Background(), []mongo.IndexModel{
			expiredIndex,
			statusIndex,
			nameIndex,
		})

	return
}

func (db *mongoJobDB) ListExpiredSchedules() (js []*Job, err error) {
	ctx := context.TODO()
	f := bson.M{
		"next_run_at": bson.M{"$lte": now()},
		"status":      PENDING.String(),
	}

	cursor, err := db.conn.
		Database(db.dbName).
		Collection(db.collName).
		Find(ctx, f)
	if err != nil {
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var jd jobDocument
		err = cursor.Decode(&jd)
		if err != nil {
			continue
		}
		j := unmarshalJob(jd)
		js = append(js, &j)
	}

	return
}

func (db *mongoJobDB) List(f Finder) (js []*Job, err error) {
	ctx := context.TODO()
	filter := parseListFilter(f)

	cursor, err := db.conn.
		Database(db.dbName).
		Collection(db.collName).
		Find(ctx, filter)
	if err != nil {
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var jd jobDocument
		err = cursor.Decode(&jd)
		if err != nil {
			continue
		}
		j := unmarshalJob(jd)
		js = append(js, &j)
	}

	return
}

func (db *mongoJobDB) SaveJob(j Job) (err error) {
	doc := marshalJob(j)

	if doc.ID == nil {
		// no ID, then its a new job
		_, err = db.conn.
			Database(db.dbName).
			Collection(db.collName).
			InsertOne(
				context.TODO(),
				doc,
			)
		return
	}

	options := options.Update().SetUpsert(true)
	update := bson.M{"$set": doc}
	filter := bson.M{"_id": doc.ID}

	_, err = db.conn.
		Database(db.dbName).
		Collection(db.collName).
		UpdateOne(
			context.TODO(),
			filter,
			update,
			options,
		)

	return
}

func (db *mongoJobDB) DeleteJob(j Job) (err error) {
	doc := marshalJob(j)
	filter := bson.M{"_id": doc.ID}

	_, err = db.conn.
		Database(db.dbName).
		Collection(db.collName).
		DeleteOne(
			context.TODO(),
			filter,
		)

	return
}

func parseListFilter(f Finder) (filter bson.M) {
	filter = bson.M{}

	if f.Name != "" {
		filter["name"] = f.Name
	}

	if f.Status != "" {
		filter["status"] = f.Status
	}

	if f.Data != nil {
		for field, value := range f.Data {
			filterField := fmt.Sprintf("data.%s", field)
			filter[filterField] = value
		}
	}

	return
}
