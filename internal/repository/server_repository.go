package repository

import (
	"context"
	"errors"
	"fmt"

	"singularity/internal/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerRepository struct {
	collection *mongo.Collection
}

func NewServerRepository(database *mongo.Database) *ServerRepository {
	return &ServerRepository{
		database.Collection("servers"),
	}
}

func (serverRepository *ServerRepository) GetAll(ctx context.Context) ([]*data.Server, error) {
	cursor, err := serverRepository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println("Error closing cursor:", err)
		}
	}(cursor, ctx)

	var servers []*data.Server
	if err := cursor.All(ctx, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func (serverRepository *ServerRepository) Insert(ctx context.Context, server *data.Server) error {
	_, err := serverRepository.collection.InsertOne(ctx, server)
	return err
}

func (serverRepository *ServerRepository) DeleteByID(ctx context.Context, id string) error {
	res, err := serverRepository.collection.DeleteOne(ctx, bson.M{"server_id": id})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("server not found")
	}

	return nil
}
