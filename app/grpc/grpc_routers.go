package grpc_routers

import (
	"context"

	"github.com/alirezakargar1380/agar.io-golang/test_log"
)

type Serverr struct {
}

func (s *Serverr) SetMessage(ctx context.Context, message *test_log.Log) (*test_log.Log, error) {
	return &test_log.Log{
		Name: "fucked up",
	}, nil
}
