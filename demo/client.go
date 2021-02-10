package main

import (
	"context"
	"log"
	"os"
	"strconv"

	pm "bitbucket.org/atlant-io/genproto/gen/go/product-manager/v1"
	tw "github.com/olekukonko/tablewriter"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := pm.NewProductManagerServiceClient(conn)

	if _, err := c.Fetch(context.Background(), &pm.FetchRequest{
		Url: "",
	}); err != nil {
		log.Println(err)
	}

	resp, err := c.List(context.Background(), &pm.ListRequest{
		NextPageToken: "",
		PageSize:      3,
		OrderBy:       "price desc,name",
	})
	if err != nil {
		log.Println(err)
	}

	table := tw.NewWriter(os.Stdout)
	table.SetHeader([]string{"№", "Product name", "Price", "Price change count", "Update unix"})

	for i, v := range resp.GetProducts() {
		table.Append([]string{
			strconv.Itoa(i),
			v.GetName(),
			v.GetPrice(),
			strconv.Itoa(int(v.GetPriceChangeCount())),
			strconv.Itoa(int(v.GetUpdateUnix())),
		})
	}

	table.Render()
}