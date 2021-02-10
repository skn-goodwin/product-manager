package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	pm "bitbucket.org/atlant-io/genproto/gen/go/product-manager/v1"
	"bitbucket.org/atlant-io/product-manager/repo"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ProductManagerServer struct {
	repo repo.Repo

	// If true: fetch local CSV file (product.csv)
	isLocal bool
}

func NewProductManagerServer(repo repo.Repo, isLocal bool) *ProductManagerServer {
	return &ProductManagerServer{
		repo:    repo,
		isLocal: isLocal,
	}
}

func StartServer(serverAddr, dbUri string, isLocal bool) {
	mc, err := repo.NewMongoClient(dbUri)
	if err != nil {
		panic(err)
	}

	gs := grpc.NewServer()

	s := NewProductManagerServer(repo.NewRepo(mc), isLocal)
	pm.RegisterProductManagerServiceServer(gs, s)

	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		panic(err)
	}

	log.Println("gRPC server started...")

	if err := gs.Serve(lis); err != nil {
		panic(err)
	}
}

func StartGateway(gatewayAddr, serverAddr string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(
		hideHeaders(),
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := pm.RegisterProductManagerServiceHandlerFromEndpoint(ctx, mux, serverAddr, opts); err != nil {
		panic(err)
	}

	log.Println("gRPC gateway started...")

	if err := http.ListenAndServe(gatewayAddr, mux); err != nil {
		panic(err)
	}
}

func hideHeaders() runtime.ServeMuxOption {
	return runtime.WithOutgoingHeaderMatcher(func(s string) (string, bool) {
		if strings.HasPrefix(strings.ToLower(s), "grpc") {
			return "", false
		}
		return s, true
	})
}