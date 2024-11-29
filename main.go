package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

const minSecondsDiff = 30

func main() {
	conn := client.NewGrpcClient("grpc.nile.trongrid.io:50051")
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		fmt.Println("error connecting", err)
		return
	}

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createTable(db)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		http.HandleFunc("/createOrder", generateAddressApi(db))
		http.HandleFunc("/getAllData", getAllAddressApi(db))
		http.HandleFunc("/transferToken", transferTokenFromMerchant(db, conn))
		http.HandleFunc("/activateAccount", activateAccount(db, conn))
		http.HandleFunc("/clientToMerchant", clientToMerchant(db, conn))
		http.HandleFunc("/getAllAccountBalance", getAllAccountBalance(db, conn))
		fmt.Println("HTTP Server running on :8000")
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	// Start the for loops in parallel
	// wg.Add(2) // Add 3 to the WaitGroup counter
	// go infinteLoopFirst(db, conn, &wg)
	// go infinteLoopSecond(db, conn, &wg)
	wg.Wait()
}

func infinteLoopFirst(db *sql.DB, conn *client.GrpcClient, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i >= 0; i++ {
		users, err := getTableData(db)
		if err != nil {
			fmt.Println("error getting table data", err)
			return
		}
		for i := 0; i < len(users); i++ {
			fmt.Println("---------------")
			fmt.Println("---------------")
			fmt.Println("---------------")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			TokenTransfer(db, conn, users[i].AddressKey, users[i].Contract, users[i].ReceivingAddress, users[i].PrivateKey, users[i].ReceivingPrivate, users[i].Id, "first")

			fmt.Println("---------------")
			fmt.Println("---------------")
			fmt.Println("---------------")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
		}
	}
}

func infinteLoopSecond(db *sql.DB, conn *client.GrpcClient, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(minSecondsDiff * time.Second)
	for i := 0; i >= 0; i++ {
		users, err := getTableData(db)
		if err != nil {
			fmt.Println("error getting table data", err)
			return
		}
		for i := 0; i < len(users); i++ {
			fmt.Println("==============")
			fmt.Println("==============")
			fmt.Println("==============")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			TokenTransfer(db, conn, users[i].AddressKey, users[i].Contract, users[i].ReceivingAddress, users[i].PrivateKey, users[i].ReceivingPrivate, users[i].Id, "second")
			fmt.Println("==============")
			fmt.Println("==============")
			fmt.Println("==============")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
		}
	}
}
