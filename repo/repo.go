package repo

import (
	"context"
	"log"
	"time"

	pm "bitbucket.org/atlant-io/genproto/gen/go/product-manager/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repo interface {
	GetClient() *mongo.Client

	UpsertProduct(product *pm.Product) error
	ListProducts(nextPageToken string, pageSize int32, orderBy string) ([]*pm.Product, error) // TODO: use Meta message ..?
}

type repo struct {
	client *mongo.Client
}

func NewRepo(client *mongo.Client) Repo {
	return &repo{
		client: client,
	}
}

func NewMongoClient(uri string) (*mongo.Client, error) {
	opts := options.Client().
		ApplyURI(uri).
		SetAuth(options.Credential{
			Username: "usr",
			Password: "pwd",
		})

	client, err := mongo.NewClient(opts)
	if err != nil {
		log.Println("[ERROR] mongo.NewClient() --> ", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Println("[ERROR] client.Connect() --> ", err)
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("[ERROR] client.Ping() --> ", err)
		return nil, err
	}

	return client, nil
}

func (r *repo) GetClient() *mongo.Client {
	return r.client
}