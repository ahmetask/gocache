package gocache

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
)

var cache *Cache

type ServerConfig struct {
	Port     string
	CachePtr *Cache
}

type Server struct {
}

func (s *Server) GetCache(ctx context.Context, request *GetCacheRequest) (*GetCacheResponse, error) {
	cache, success := cache.Get(request.Key)
	if !success {
		return nil, NewError(404, fmt.Sprintf("can not get cache with key:%s", request.Key))
	}

	response := &GetCacheResponse{
		Value:   cache.([]byte),
		Success: success,
	}

	return response, nil
}

func (s *Server) SaveCache(ctx context.Context, request *SaveCacheRequest) (*Empty, error) {
	success, message := cache.Add(request.Key, request.Value, time.Second*time.Duration(request.Life))
	if !success {
		return nil, NewError(409, fmt.Sprintf("can not save e:%s", message))
	}

	return &Empty{}, nil
}

func (s *Server) DeleteCachedItem(ctx context.Context, request *DeleteCacheItemRequest) (*Empty, error) {
	cache.DeleteCachedItem(request.Key)
	return &Empty{}, nil
}

func (s *Server) ClearCache(ctx context.Context, request *Empty) (*Empty, error) {
	cache.ClearCache()
	return &Empty{}, nil
}

func validateServerConfig(config ServerConfig) (bool, *Error) {
	if config.CachePtr == nil {
		return false, NewError(400, "you should add cache instance")
	}

	if config.Port == "" {
		return false, NewError(400, "port should not be empty")
	}

	return true, nil
}

func NewCacheServer(config ServerConfig) *Error {
	valid, validationErr := validateServerConfig(config)

	if !valid {
		return validationErr
	}

	cache = config.CachePtr

	address := "0.0.0.0:" + config.Port
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return NewError(500, err.Error())
	}

	s := grpc.NewServer()
	RegisterCacheServiceServer(s, &Server{})

	serveError := s.Serve(listen)

	if serveError != nil {
		return NewError(500, serveError.Error())
	}

	return nil
}
