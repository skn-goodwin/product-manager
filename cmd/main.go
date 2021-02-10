package main

import (
	"os"
	"strconv"

	"bitbucket.org/atlant-io/product-manager/server"
)

func main() {
	serverAddr := os.Getenv("SERVER_ADDR")
	dbUri := os.Getenv("DB_URI")
	gatewayAddr := os.Getenv("GATEWAY_ADDR")
	isLocal, _ := strconv.ParseBool(os.Getenv("IS_LOCAL"))

	go server.StartServer(serverAddr, dbUri, isLocal)
	go server.StartGateway(gatewayAddr, serverAddr)

	select {}
}