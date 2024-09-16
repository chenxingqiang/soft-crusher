// File: internal/database/mongodb/mongodb.go

package mongodb

import (
	"context"

	"github.com/chenxingqiang/soft-crusher/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// previous MongoDB struct and NewMongoDB function
type MongoDB struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoDB(connectionString, database, collection string) (*MongoDB, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (m *MongoDB) SaveSoftwareInfo(info *models.SoftwareInfo) error {
	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.InsertOne(context.Background(), info)
	return err
}

// Part 2: MongoDB methods implementation

func (m *MongoDB) SaveSoftwareInfo(info *models.SoftwareInfo) error {
	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.InsertOne(context.Background(), info)
	return err
}

func (m *MongoDB) GetSoftwareInfo(id string) (*models.SoftwareInfo, error) {
	coll := m.client.Database(m.database).Collection(m.collection)
	var info models.SoftwareInfo
	err := coll.FindOne(context.Background(), bson.M{"_id": id}).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (m *MongoDB) UpdateSoftwareInfo(info *models.SoftwareInfo) error {
	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.ReplaceOne(context.Background(), bson.M{"_id": info.ID}, info)
	return err
}

func (m *MongoDB) DeleteSoftwareInfo(id string) error {
	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (m *MongoDB) ListSoftwareInfo(limit, offset int) ([]*models.SoftwareInfo, error) {
	coll := m.client.Database(m.database).Collection(m.collection)
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := coll.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*models.SoftwareInfo
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *MongoDB) SearchSoftwareInfo(query string) ([]*models.SoftwareInfo, error) {
	coll := m.client.Database(m.database).Collection(m.collection)
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"keywords": bson.M{"$regex": query, "$options": "i"}},
		},
	}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*models.SoftwareInfo
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *MongoDB) GetSoftwareInfoByCodeRepository(url string) (*models.SoftwareInfo, error) {
	coll := m.client.Database(m.database).Collection(m.collection)
	var info models.SoftwareInfo
	err := coll.FindOne(context.Background(), bson.M{"codeRepository": url}).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (m *MongoDB) SearchSoftwareInfo(query string) ([]*models.SoftwareInfo, error) {
	coll := m.client.Database(m.database).Collection(m.collection)
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"keywords": bson.M{"$regex": query, "$options": "i"}},
		},
	}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*models.SoftwareInfo
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *MongoDB) GetSoftwareInfoByCodeRepository(url string) (*models.SoftwareInfo, error) {
	coll := m.client.Database(m.database).Collection(m.collection)
	var info models.SoftwareInfo
	err := coll.FindOne(context.Background(), bson.M{"codeRepository": url}).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
