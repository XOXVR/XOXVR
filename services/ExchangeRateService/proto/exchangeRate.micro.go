// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/exchangeRate.proto

/*
Package auth is a generated protocol buffer package.

It is generated from these files:
	proto/exchangeRate.proto

It has these top-level messages:
	GetRateReq
	GetRateResp
*/
package auth

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Currency service

type CurrencyService interface {
	GetRate(ctx context.Context, in *GetRateReq, opts ...client.CallOption) (*GetRateResp, error)
}

type currencyService struct {
	c    client.Client
	name string
}

func NewCurrencyService(name string, c client.Client) CurrencyService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "auth"
	}
	return &currencyService{
		c:    c,
		name: name,
	}
}

func (c *currencyService) GetRate(ctx context.Context, in *GetRateReq, opts ...client.CallOption) (*GetRateResp, error) {
	req := c.c.NewRequest(c.name, "Currency.GetRate", in)
	out := new(GetRateResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Currency service

type CurrencyHandler interface {
	GetRate(context.Context, *GetRateReq, *GetRateResp) error
}

func RegisterCurrencyHandler(s server.Server, hdlr CurrencyHandler, opts ...server.HandlerOption) {
	type currency interface {
		GetRate(ctx context.Context, in *GetRateReq, out *GetRateResp) error
	}
	type Currency struct {
		currency
	}
	h := &currencyHandler{hdlr}
	s.Handle(s.NewHandler(&Currency{h}, opts...))
}

type currencyHandler struct {
	CurrencyHandler
}

func (h *currencyHandler) GetRate(ctx context.Context, in *GetRateReq, out *GetRateResp) error {
	return h.CurrencyHandler.GetRate(ctx, in, out)
}
