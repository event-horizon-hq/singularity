package repository

import (
	"context"
	"errors"
	"fmt"

	"singularity/internal/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServerRepository struct {
	collection *mongo.Collection
}

func NewServerRepository(database *mongo.Database) *ServerRepository {
	return &ServerRepository{
		database.Collection("servers"),
	}
}

func (repository *ServerRepository) EnsureIndexes(ctx context.Context) error {
	_, err := repository.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "discriminator", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetName("unique_discriminator"),
	})

	return err
}

func (repository *ServerRepository) GetAll(ctx context.Context) ([]*data.Server, error) {
	cursor, err := repository.collection.Find(ctx, bson.M{})
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

func (repository *ServerRepository) Insert(ctx context.Context, server *data.Server) error {
	filter := bson.M{"discriminator": server.Discriminator}

	update := bson.M{
		"$set": server,
	}

	opts := options.Update().SetUpsert(true)

	_, err := repository.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (repository *ServerRepository) DeleteByID(ctx context.Context, id string) error {
	res, err := repository.collection.DeleteOne(ctx, bson.M{"discriminator": id})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("server not found")
	}

	return nil
}
