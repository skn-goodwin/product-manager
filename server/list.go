package server

import (
	"context"
	"log"

	pm "bitbucket.org/atlant-io/genproto/gen/go/product-manager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ProductManagerServer) List(ctx context.Context, r *pm.ListRequest) (*pm.ListResponse, error) {
	products, err := s.repo.ListProducts(r.GetNextPageToken(), r.GetPageSize(), r.GetOrderBy())
	if err != nil {
		log.Println("[ERROR] s.repo.ListProducts() --> ", err)
		return nil, status.Error(codes.Unknown, "Oops..")
	}

	return &pm.ListResponse{
		Products: products,
	}, nil
}