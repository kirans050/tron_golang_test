package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/mr-tron/base58"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/mattn/go-sqlite3"
)

type keysStruct struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Address    string `json:"address"`
}

type AddressTable struct {
	Id               int      `json:"id"`
	PublicKey        string   `json:"publicKey"`
	PrivateKey       string   `json:"privateKey"`
	AddressKey       string   `json:"address"`
	Amount           float64  `json:"amount"`
	TimeStamp        int64    `json:"timestamp"`
	Token            *string  `json:"token"`
	OrderID          *float64 `json:"order_id"`
	CallBack         string   `json:"callback"`
	ReceivingAddress string   `json:"receiving_address"`
	Contract         string   `json:"contract"`
	TrxTimeStamp     int64    `json:"trxTimeStamp"`
	ReceivingPrivate string   `json:"receiving_private"`
	CallbackCount    int64    `json:"callBackCount"`
}

type OrderRequest struct {
	OrderID          int    `json:"order_id"`
	Token            string `json:"token"`
	Amount           int    `json:"amount"`
	Callback         string `json:"callback"`
	ReceivingAddress string `json:"receiving_address"`
	Contract         string `json:"contract"`
	ReceivingPrivate string `json:"receiving_private"`
	CallbackCount    int    `json:"callbackcount"`
}

func generateAddressHandler() (keysStruct, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
		return keysStruct{}, err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateAddress := hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
		return keysStruct{}, err
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicAddress := hexutil.Encode(publicKeyBytes)[2:]

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	addressHex := "41" + address[2:]
	addb, _ := hex.DecodeString(addressHex)
	hash1 := s256(s256(addb))
	secret := hash1[:4]
	for _, v := range secret {
		addb = append(addb, v)
	}
	addressBase58 := base58.Encode(addb)
	keys := keysStruct{
		PrivateKey: privateAddress,
		PublicKey:  publicAddress,
		Address:    addressBase58,
	}
	return keys, nil
}

func insertTableData(db *sql.DB, keys keysStruct, order OrderRequest) (int64, error) {
	//  Insert some data into the table
	stmt, err := db.Prepare("INSERT INTO addresses(publicKey, privateKey, address, amount, timestamp, token, order_id, callback, receiving_address, contract, receiving_private, callBackCount) VALUES(?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer stmt.Close()

	//  Insert data
	result, err := stmt.Exec(keys.PublicKey, keys.PrivateKey, keys.Address, order.Amount, time.Now().Unix(), order.Token, order.OrderID, order.Callback, order.ReceivingAddress, order.Contract, order.ReceivingPrivate, order.CallbackCount)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	insertedID, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Failed to get last insert ID:", err)
		return 0, err
	}
	return insertedID, nil
}

func getTableData(db *sql.DB) ([]AddressTable, error) {
	//  Query the data and print it
	rows, err := db.Query("SELECT * FROM addresses")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	//  Print the results
	var AddressSlice []AddressTable
	for rows.Next() {
		var Address AddressTable
		err := rows.Scan(
			&Address.Id,
			&Address.PublicKey,
			&Address.PrivateKey,
			&Address.AddressKey,
			&Address.Amount,
			&Address.TimeStamp,
			&Address.Token,
			&Address.OrderID,
			&Address.CallBack,
			&Address.ReceivingAddress,
			&Address.Contract,
			&Address.TrxTimeStamp,
			&Address.ReceivingPrivate,
			&Address.CallbackCount)
		if err != nil {
			return nil, err
		}
		AddressSlice = append(AddressSlice, Address)
	}

	//  Check for errors in the row iteration
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return AddressSlice, nil
}

func createTable(db *sql.DB) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS addresses (
		id INTEGER PRIMARY KEY,
		publicKey TEXT,
 		privateKey TEXT,
		address TEXT,
		amount INTEGER DEFAULT 0.0,
		timestamp INTEGER DEFAULT (strftime('%s', 'now')),
		token TEXT DEFAULT "",
		order_id INTEGER ,
		callback TEXT TEXT DEFAULT "",
		receiving_address TEXT DEFAULT "",
		contract TEXT DEFAULT "",
		trxTimeStamp INTEGER DEFAULT 0,
		receiving_private TEXT DEFAULT "",
		callBackCount INTEGER DEFAULT 0
 	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

func s256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "hello world")
}

func generateAddressApi(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var order OrderRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&order)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		keys, err := generateAddressHandler()
		if err != nil {
			http.Error(w, "failed to generate keys", http.StatusInternalServerError)
			return
		}
		var row AddressTable
		insertedID, err := insertTableData(db, keys, order)
		if err != nil {
			http.Error(w, "failed to create data", http.StatusInternalServerError)
			return
		}

		query := `SELECT id, publicKey, privateKey, address, amount, timestamp, token, order_id, callback, receiving_address, contract, receiving_private, callBackCount FROM addresses WHERE id = ?`
		err = db.QueryRow(query, insertedID).Scan(&row.Id,
			&row.PublicKey, &row.PrivateKey, &row.AddressKey, &row.Amount, &row.TimeStamp,
			&row.Token, &row.OrderID, &row.CallBack, &row.ReceivingAddress, &row.Contract,
			&row.ReceivingPrivate, &row.CallbackCount,
		)
		if err != nil {
			http.Error(w, "Failed to query inserted row", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(row)
	}
}

func getAllAddressApi(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  Fetch users from the database
		users, err := getTableData(db)
		if err != nil {
			log.Fatal("error", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		//  Convert the users slice to JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func transferTokenFromMerchant(db *sql.DB, conn *client.GrpcClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  Fetch users from the database
		users, err := getTableData(db)
		if err != nil {
			log.Fatal("error", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		contract := "TY1DBj7Ys1bDcK37kwATaQpHxdTCnYrr1f"
		merchantAccPrivate := "17c112793ba29f39dc0b6056695746a76f19bd8eb1e695d88d3c2dfdb30edb42"
		merchantAccAddress := "TWYywngN3EfYiyY2NHzAHi4ad9B1uJNb8Y"
		for i := 0; i < len(users); i++ {
			merchantToClientToken(conn, users[i].AddressKey, contract, merchantAccAddress, users[i].PrivateKey, merchantAccPrivate)
		}

		//  Convert the users slice to JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
func activateAccount(db *sql.DB, conn *client.GrpcClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  Fetch users from the database
		users, err := getTableData(db)
		if err != nil {
			log.Fatal("error", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		merchantAccPrivate := "17c112793ba29f39dc0b6056695746a76f19bd8eb1e695d88d3c2dfdb30edb42"
		merchantAccAddress := "TWYywngN3EfYiyY2NHzAHi4ad9B1uJNb8Y"
		for i := 0; i < len(users); i++ {
			ActivateNewAccount(conn, users[i].AddressKey, merchantAccAddress, merchantAccPrivate)
		}

		//  Convert the users slice to JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
func clientToMerchant(db *sql.DB, conn *client.GrpcClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  Fetch users from the database
		users, err := getTableData(db)
		if err != nil {
			log.Fatal("error", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		for i := 0; i < len(users); i++ {
			TokenTransfer(db, conn, users[i].AddressKey, users[i].Contract, users[i].ReceivingAddress, users[i].PrivateKey, users[i].ReceivingPrivate, users[i].Id, "generate")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func getAllAccountBalance(db *sql.DB, conn *client.GrpcClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  Fetch users from the database
		users, err := getTableData(db)
		if err != nil {
			log.Fatal("error", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		contract := "TY1DBj7Ys1bDcK37kwATaQpHxdTCnYrr1f"
		for i := 0; i < len(users); i++ {

			var trc20Balance = big.NewInt(0)
			var trxBalance float64 = 0
			token, err := conn.TRC20ContractBalance(users[i].AddressKey, contract)
			if err != nil {
				fmt.Println("error getting token balance", users[i].Id, users[i].AddressKey)
				token = trc20Balance
			}
			trx, err := GetAccountBalance(conn, users[i].AddressKey)
			if err != nil {
				fmt.Println("error getting trx balance", users[i].Id, users[i].AddressKey)
				trx = trxBalance
			}

			fmt.Println("balance", users[i].AddressKey, token, "--", trx)
		}

		//  Convert the users slice to JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
