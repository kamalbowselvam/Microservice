package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	protos "GitHub/Microservice/currency/protos/currency"
)


type Currency struct {
	log hclog.Logger
}


func (c* Currency) GetRate(ctx context.Context, req* protos.RateRequest)(*protos.RateResponse, error)