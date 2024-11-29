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
	// infinteLoop(db, conn)
	log.Println("Server running at http:localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func infinteLoop(db *sql.DB, conn *client.GrpcClient) {
	for i := 0; i >= 0; i++ {
		users, err := getTableData(db)
		if err != nil {
			fmt.Println("error getting table data", err)
			return
		}
		contract := "TY1DBj7Ys1bDcK37kwATaQpHxdTCnYrr1f"
		merchantAccPrivate := "17c112793ba29f39dc0b6056695746a76f19bd8eb1e695d88d3c2dfdb30edb42"
		merchantAccAddress := "TWYywngN3EfYiyY2NHzAHi4ad9B1uJNb8Y"
		for i := 0; i < len(users); i++ {
			TokenTransfer(db, conn, users[i].AddressKey, contract, merchantAccAddress, users[i].PrivateKey, merchantAccPrivate, users[i].Id)
		}
	}
}

// 3.09411(0.345 + 2.74911) + 0.268
// 3.36211

// 2.74911 + 0.268
// 3.01711
// 4.1711
