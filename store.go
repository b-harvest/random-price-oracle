package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoreService struct {
	cfg MongoDBConfig
	mc  *mongo.Client
}

func NewStoreService(cfg MongoDBConfig, mc *mongo.Client) *StoreService {
	return &StoreService{cfg: cfg, mc: mc}
}

func (s *StoreService) CoinCollection() *mongo.Collection {
	return s.mc.Database(s.cfg.DB).Collection(s.cfg.CoinCollection)
}

func (s *StoreService) CreateIndexes(ctx context.Context) error {
	_, err := s.CoinCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{CoinSymbolKey, 1}},
	})
	return err
}

func (s *StoreService) Coins(ctx context.Context) ([]Coin, error) {
	cur, err := s.CoinCollection().Find(ctx, bson.M{
		CoinSymbolKey: bson.M{"$exists": true},
		CoinPriceKey:  bson.M{"$gt": 0.0},
	})
	if err != nil {
		return nil, err
	}
	var coins []Coin
	if err := cur.All(ctx, &coins); err != nil {
		return nil, err
	}
	return coins, nil
}

func (s *StoreService) Coin(ctx context.Context, symbol string) (Coin, error) {
	symbol = NormalizeSymbol(symbol)
	var c Coin
	if err := s.CoinCollection().FindOne(ctx, bson.M{
		CoinSymbolKey: symbol,
	}).Decode(&c); err != nil {
		return Coin{}, err
	}
	return c, nil
}

func (s *StoreService) SetPrice(ctx context.Context, symbol string, price float64) error {
	symbol = NormalizeSymbol(symbol)
	_, err := s.CoinCollection().ReplaceOne(ctx, bson.M{
		CoinSymbolKey: symbol,
	}, Coin{
		Symbol:    symbol,
		Price:     price,
		UpdatedAt: time.Now(),
	}, options.Replace().SetUpsert(true))
	return err
}
