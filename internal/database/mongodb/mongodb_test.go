package mongodb

import (
	"context"
	"testing"

	"github.com/chenxingqiang/soft-crusher/internal/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDB(t *testing.T) {
	// Connect to a test MongoDB instance
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)
	defer client.Disconnect(context.Background())

	db := NewMongoDB(client.Database("test_db"))

	// Test SaveSoftwareInfo
	info := &models.SoftwareInfo{Name: "Test Software"}
	err = db.SaveSoftwareInfo(info)
	assert.NoError(t, err)
	assert.NotEmpty(t, info.ID)

	// Test GetSoftwareInfo
	retrievedInfo, err := db.GetSoftwareInfo(info.ID)
	assert.NoError(t, err)
	assert.Equal(t, info.Name, retrievedInfo.Name)

	// Test UpdateSoftwareInfo
	info.Name = "Updated Software"
	err = db.UpdateSoftwareInfo(info)
	assert.NoError(t, err)

	retrievedInfo, err = db.GetSoftwareInfo(info.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Software", retrievedInfo.Name)

	// Test DeleteSoftwareInfo
	err = db.DeleteSoftwareInfo(info.ID)
	assert.NoError(t, err)

	_, err = db.GetSoftwareInfo(info.ID)
	assert.Error(t, err)

	// Add more tests for other database operations...
}
