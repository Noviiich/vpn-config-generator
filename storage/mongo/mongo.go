package mongo

import (
	"context"
	"log"
	"time"

	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	configs Configs
}

type Configs struct {
	Collection *mongo.Collection
}

type Config struct {
	ID         string
	PrivateKey string
	PublicKey  string
	IPAddress  string
}

func New(connectString string, connectTimeout time.Duration) Storage {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))
	if err != nil {
		log.Fatalf("can't connect mongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("can't ping mongoDB: %v", err)
	}

	configs := Configs{
		Collection: client.Database("vpn").Collection("configs"),
	}

	return Storage{
		configs: configs,
	}
}

func (s Storage) Create(ctx context.Context, conf *storage.Config) error {
	_, err := s.configs.Collection.InsertOne(ctx, conf)
	if err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}
