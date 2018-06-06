//go:generate abigen --sol ./contracts/Crowdsale.sol --pkg gocontracts --out ./gocontracts/Crowdsale.go

// go:generate protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. ./services/FiatService/proto/fiatService.proto
// go:generate protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. services/FiatService/proto/fiatService.proto

package main
