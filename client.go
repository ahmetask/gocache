package gocache

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
)

type AddCacheRequest struct {
	Key   string
	Value interface{}
	Life  int64 //as second
}

func GetCache(address, key string, cacheResponse interface{}) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	client := NewCacheServiceClient(cc)
	request := &GetCacheRequest{Key: key}

	resp, clientError := client.GetCache(context.Background(), request)
	if clientError != nil {
		return clientError
	}

	marshallError := json.Unmarshal(resp.Value, &cacheResponse)
	if marshallError != nil {
		return marshallError
	}

	return nil
}

func SaveCache(address string, cache AddCacheRequest) error {

	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	client := NewCacheServiceClient(cc)

	cacheBytes, bodyError := json.Marshal(&cache.Value)

	if bodyError != nil {
		return bodyError
	}

	request := &SaveCacheRequest{Key: cache.Key, Value: cacheBytes, Life: cache.Life}

	_, clientError := client.SaveCache(context.Background(), request)
	if clientError != nil {
		return clientError
	}

	return nil
}

func ClearCache(address string) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	client := NewCacheServiceClient(cc)
	request := &Empty{}

	_, clientError := client.ClearCache(context.Background(), request)
	if clientError != nil {
		return clientError
	}

	return nil
}

func DeleteCacheItem(address string, key string) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	client := NewCacheServiceClient(cc)
	request := &DeleteCacheItemRequest{Key: key}

	_, clientError := client.DeleteCachedItem(context.Background(), request)
	if clientError != nil {
		return clientError
	}

	return nil
}
