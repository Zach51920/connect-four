package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type Config struct {
	URI    string
	DBName string
}

type Provider struct {
	config *Config
	client *mongo.Client
}

func FromEnv() *Config {
	return &Config{
		URI:    os.Getenv("MONGO_URI"),
		DBName: os.Getenv("MONGO_DB_NAME"),
	}
}

func NewProvider(config *Config) (*Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(config.URI).
		SetServerAPIOptions(serverAPI).
		SetMaxPoolSize(20).
		SetHeartbeatInterval(30 * time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	provider := &Provider{client: client, config: config}
	return provider, nil
}

func (p *Provider) Client() *mongo.Client {
	return p.client
}

func (p *Provider) DB() *mongo.Database {
	return p.client.Database(p.config.DBName)
}

func (p *Provider) Ping() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second)
	defer ctxCancel()

	var result bson.D
	return p.DB().
		RunCommand(ctx, bson.D{{"ping", 1}}).
		Decode(&result)
}

func (p *Provider) Close() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	return p.client.Disconnect(ctx)
}
