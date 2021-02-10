package server

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pm "bitbucket.org/atlant-io/genproto/gen/go/product-manager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var field2Idx = map[string]int{
	"PRODUCT NAME": -1,
	"PRICE":        -1,
}

func (s *ProductManagerServer) Fetch(ctx context.Context, r *pm.FetchRequest) (*emptypb.Empty, error) {
	now := time.Now().Unix()

	var ioReader io.Reader
	if s.isLocal {
		f, err := os.Open("product.csv")
		if err != nil {
			log.Println("[ERROR] os.Open() --> ", err)
			return nil, status.Error(codes.Unknown, "Oops..")
		}

		ioReader = f
	} else {
		// TODO: validate url (prevent access to the internal network)

		resp, err := http.Get(r.GetUrl())
		if err != nil {
			log.Println("[ERROR] http.Get() --> ", err)
			return nil, status.Error(codes.Unknown, "Oops..")
		}
		defer resp.Body.Close()

		ioReader = resp.Body
	}

	reader := csv.NewReader(ioReader)
	reader.Comma = ';'

	i := 0
	productNameIdx := 0
	productPriceIdx := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("[ERROR] reader.Read() --> ", err)
			return nil, status.Error(codes.Unknown, "Oops..")
		}

		// Skip headers in CSV file
		// Get indexes for each field
		if i == 0 {
			FillFieldsIndex(record)
			productNameIdx  = field2Idx["PRODUCT NAME"]
			productPriceIdx = field2Idx["PRICE"]
			i = 1
			continue
		}

		if err := s.repo.UpsertProduct(&pm.Product{
			Name:       record[productNameIdx],
			Price:      record[productPriceIdx],
			UpdateUnix: now,
		}); err != nil {
			log.Println("[ERROR] s.repo.UpsertProduct() --> ", err)
			return nil, status.Error(codes.Unknown, "Oops..")
		}
	}

	return new(emptypb.Empty), nil
}

func FillFieldsIndex(csvHeader []string) {
	for i, header := range csvHeader {
		field2Idx[header] = i
	}
}