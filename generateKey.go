package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	Id               int     `json:"id"`
	PublicKey        string  `json:"publicKey"`
	PrivateKey       string  `json:"privateKey"`
	AddressKey       string  `json:"address"`
	Amount           float64 `json:"amount"`
	TimeStamp        int     `json:"timestamp"`
	Token            *string `json:"token"`
	CallbacksCount   float64 `json:"callbacks_count"`
	OrderId          float64 `json:"order_id"`
	CallBack         string  `json:"callback"`
	ReceivingAddress string  `json:"reciving_address"`
	Contract         string  `json:"contract"`
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

func insertTableData(db *sql.DB, keys keysStruct) (sql.Result, error) {
	//  Insert some data into the table
	stmt, err := db.Prepare("INSERT INTO addresses(publicKey, privateKey, address) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	//  Insert data
	result, err := stmt.Exec(keys.PublicKey, keys.PrivateKey, keys.Address)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result, nil
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
		err := rows.Scan(&Address.Id, &Address.PublicKey, &Address.PrivateKey, &Address.AddressKey, &Address.Amount, &Address.TimeStamp, &Address.Token, &Address.CallbacksCount, &Address.OrderId, &Address.CallBack, &Address.ReceivingAddress, &Address.Contract)
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
		callback TEXT TEXT DEFAULT "",
		order_id INTEGER ,
		reciving_address TEXT DEFAULT "",
		contract TEXT DEFAULT ""

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
		keys, err := generateAddressHandler()
		if err != nil {
			http.Error(w, "failed to generate keys", http.StatusInternalServerError)
			return
		}
		_, err = insertTableData(db, keys)
		if err != nil {
			http.Error(w, "failed to insert data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keys)
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
