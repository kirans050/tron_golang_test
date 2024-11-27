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

	// fromAddress := "TVJfeB4se3BHzL5Vr6Fya94K3MoWpf6WRU"
	// // // // // // address := "TXLL89KocCLcoDM7Co38LSqAqwHUzT9VuM"
	// privateKey := "428f90cc7476f408280379b9349ea1d514b244ea92c57bdc753deafb7e072fd8"
	// privateKey := "fd833288ba9fc8817c0b086109a9a63bf9e06f419166de275876914c85095ed5"

	contract := "TY1DBj7Ys1bDcK37kwATaQpHxdTCnYrr1f"
	clientAccPrivate := "125da0f330a2ce0410f0dffac23c2956aedeb2c1ec74b1cb6bce6decb6f0e704"
	clientAccAddress := "TFSaboWQCm7XC9Cjznn1fnLXVY4tLFKXBw"

	merchantAccPrivate := "17c112793ba29f39dc0b6056695746a76f19bd8eb1e695d88d3c2dfdb30edb42"
	merchantAccAddress := "TWYywngN3EfYiyY2NHzAHi4ad9B1uJNb8Y"

	// {
	// 	"privateKey": "17c112793ba29f39dc0b6056695746a76f19bd8eb1e695d88d3c2dfdb30edb42",
	// 	"publicKey": "04a8b4acc5831278dc9f477715c48448a3dd14f38481b0c534f668b9a0923f430feea3832a5c8da5578b95ac15459a5ffe3c8095437b5522691f66e9a404161b56",
	// 	"address": "TWYywngN3EfYiyY2NHzAHi4ad9B1uJNb8Y"
	//   }

	// SendTrx(conn, toAddress, fromAddress, privateKey)
	TokenTransfer(conn, clientAccAddress, contract, merchantAccAddress, clientAccPrivate, merchantAccPrivate)

	// EstimateTransactionEnergy(conn, fromAddress, contract, toAddress)
	// serverFunction()

	// ActivateNewAccount(conn, clientAccAddress, merchantAccAddress, merchantAccPrivate)
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
