package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

func main() {
	conn := client.NewGrpcClient("grpc.nile.trongrid.io:50051")
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		fmt.Println("error connecting", err)
		return
	}
	// toAddress := "TWYywngN3EfYiyY2NHzAHi4ad9B1uJNb8Y"

	// fromAddress := "TFXGCn4QF5zbshzHUTkK9j6yKyDsQbdukZ"
	// // // // // address := "TXLL89KocCLcoDM7Co38LSqAqwHUzT9VuM"
	// contract := "TY1DBj7Ys1bDcK37kwATaQpHxdTCnYrr1f"
	// privateKey := "1ecbf87330d1b636811c4eac678f52813bea326818a90d3de226015f51ad3d34"
	// privateKey := "fd833288ba9fc8817c0b086109a9a63bf9e06f419166de275876914c85095ed5"

	// {
	// 	"privateKey": "17c112793ba29f39dc0b6056695746a76f19bd8eb1e695d88d3c2dfdb30edb42",
	// 	"publicKey": "04a8b4acc5831278dc9f477715c48448a3dd14f38481b0c534f668b9a0923f430feea3832a5c8da5578b95ac15459a5ffe3c8095437b5522691f66e9a404161b56",
	// 	"address": "TWYywngN3EfYiyY2NHzAHi4ad9B1uJNb8Y"
	//   }

	// SendTrx(conn, toAddress, fromAddress, privateKey)
	// TokenTransfer(conn, fromAddress, contract, toAddress, privateKey)

	// EstimateTransactionEnergy(conn, fromAddress, contract, toAddress)
	// serverFunction()

	// ActivateNewAccount(conn, fromAddress, toAddress, privateKey)
}

func serverFunction() {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createTable(db)
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/createOrder", generateAddressApi(db))
	http.HandleFunc("/getAllData", getAllAddressApi(db))
	log.Println("Server running at http:localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
