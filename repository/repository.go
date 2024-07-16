package repository

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"template/model"
	"template/database"
)

type Repository struct {
	collection *mongo.Collection
}

type RepositoryInterface interface {
	Create(ctx context.Context, data *model.Data) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, updatedData *model.Data) error
	Get(ctx context.Context, id string) (*model.Data, error)
	GetAll(ctx context.Context) ([]*model.Data, error)
}

func NewRepository(client *mongo.Client, dbName, collectionName string) RepositoryInterface {
	collection := client.Database(dbName).Collection(collectionName)
	return &Repository{collection: collection}
}

func (r *Repository) Create(ctx context.Context, data *model.Data) error {
	_, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("could not create: %v", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("could not delete: %v", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, id string, updatedData *model.Data) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedData}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("could not update data: %v", err)
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Data, error) {
	var result model.Data
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("could not find data: %v", err)
	}

	return &result, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*model.Data, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("could not find data: %v", err)
	}
	defer cursor.Close(ctx)

	var results []*model.Data
	for cursor.Next(ctx) {
		var data model.Data
		if err := cursor.Decode(&data); err != nil {
			return nil, fmt.Errorf("error decoding data: %v", err)
		}
		results = append(results, &data)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return results, nil
}

func Update(c *gin.Context) {
	client := database.GetClient()
	repo := NewRepository(client, "templateDb", "templateCollection")

	id := c.Param("id")

	var updatedData model.Data
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repo.Update(context.Background(), id, &updatedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update data: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})
}
