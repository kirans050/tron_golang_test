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
	// serverFunction(conn)

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	serverFunction(conn)
}

func serverFunction(conn *client.GrpcClient) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createTable(db)
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/createOrder", generateAddressApi(db))
	http.HandleFunc("/getAllData", getAllAddressApi(db))
	http.HandleFunc("/transferToken", transferTokenFromMerchant(db, conn))
	http.HandleFunc("/activateAccount", activateAccount(db, conn))
	http.HandleFunc("/clientToMerchant", clientToMerchant(db, conn))
	http.HandleFunc("/getAllAccountBalance", getAllAccountBalance(db, conn))
	log.Println("Server running at http:localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
