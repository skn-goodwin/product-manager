package repo

import (
	"context"
	"log"
	"strconv"
	"strings"

	pm "bitbucket.org/atlant-io/genproto/gen/go/product-manager/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *repo) UpsertProduct(product *pm.Product) error {
	price, err := strconv.ParseFloat(product.GetPrice(), 64)
	if err != nil {
		return err
	}

	collection := r.client.
		Database("db"). // TODO: r.GetDatabaseName() ..?
		Collection("product")

	filter := bson.D{{"name", product.GetName()}}
	update := bson.D{
		{"$inc", bson.D{
			{"price_change_count", 1},
		}},
		{"$set", bson.D{
			{"price", price},
			{"update_unix", product.GetUpdateUnix()},
		}},
	}
	opts := options.Update().SetUpsert(true)

	if _, err := collection.UpdateOne(context.Background(), filter, update, opts); err != nil {
		log.Println("[ERROR] collection.UpdateOne() --> ", err)
		return err
	}

	return nil
}

func (r *repo) ListProducts(nextPageToken string, pageSize int32, orderBy string) ([]*pm.Product, error) {
	var filter bson.M
	if len(nextPageToken) > 0 {
		// TODO: validate nextPageToken
		filter = map[string]interface{}{
			"_id": bson.D{
				{"$lt", nextPageToken},
			},
		}
	}

	// TODO: move to separate func
	var limit int64 = 10 // TODO: set default limit as constant or pass it as flag
	if pageSize > 0 {
		limit = int64(pageSize)
	}

	// TODO: move to separate func
	m := make(map[string]interface{})
	if len(orderBy) > 0 {
		params := strings.Split(strings.Join(strings.Fields(orderBy), " "), ",")
		for _, v := range params {
			parts := strings.Split(v, " ")
			switch len(parts) {
			case 2:
				if parts[1] == "desc" {
					m[parts[0]] = -1
				} else if parts[1] == "ask" {
					m[parts[0]] = 1
				}
			case 1:
				m[parts[0]] = 1
			default:
				continue
			}
		}
	}

	// Default sorting ..?
	if len(m) == 0 {
		m["name"] = 1
	}

	collection := r.client.
		Database("db"). // TODO: r.GetDatabaseName() ..?
		Collection("product")

	opts := options.Find().
		SetSort(m).
		SetLimit(limit)

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println("[ERROR] collection.Find() --> ", err)
		return nil, err
	}

	var resp []bson.M
	if err := cursor.All(context.Background(), &resp); err != nil {
		log.Println("[ERROR] cursor.All() --> ", err)
		return nil, err
	}

	// TODO: validate response schema
	products := make([]*pm.Product, len(resp))
	for i, v := range resp {
		products[i] = &pm.Product{
			Name:             v["name"].(string),
			Price:            strconv.FormatFloat(v["price"].(float64), 'f', -1, 64),
			PriceChangeCount: v["price_change_count"].(int32) - 1, // 1 is default value (not 0)
			UpdateUnix:       v["update_unix"].(int64),
		}
	}

	return products, nil
}