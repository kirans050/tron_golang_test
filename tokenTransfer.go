package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/protobuf/proto"
)

func TokenTransfer(db *sql.DB, conn *client.GrpcClient, clientAccAddress, contract, merchantAccAddress, clientAccPrivate, merchantAccPrivate string, id int) {
	balance, err := conn.TRC20ContractBalance(clientAccAddress, contract)
	if err != nil {
		fmt.Println("error checking token balance", err)
		return
	}
	fmt.Println("balance", balance)
	if balance.Cmp(big.NewInt(0)) <= 0 {
		fmt.Println("insufficent token balance")
		return
	}

	// balance := big.NewInt(125000000)

	var userAccBalance float64 = 0
	accountActivated := true
	accBalance, err := GetAccountBalance(conn, clientAccAddress)
	if err != nil {
		fmt.Println("error getting balance", err)
		userAccBalance = 0
		accountActivated = false
	} else {
		userAccBalance = float64(accBalance)
	}
	fmt.Println("acc accountActivated", accBalance, accountActivated)

	// var bal = big.NewInt(10)
	tx, err := conn.TRC20Send(clientAccAddress, merchantAccAddress, contract, balance, 10000000)
	if err != nil {
		fmt.Println("error generating transaction", err)
		return
	}

	tx, err = generateSignature(tx, clientAccPrivate)
	if err != nil {
		fmt.Println("error generating signature", err)
		return
	}
	totalBytes, err := calculateBytes(tx)
	if err != nil {
		fmt.Println("error calculating bytes", err)
		return
	}
	fmt.Println("totolBytes", totalBytes)

	resource, err := GetAccountResourceHandler(conn, clientAccAddress)
	if err != nil {
		fmt.Println("unable to get the account resource")
		return
	}

	var totalTrxNeeded float64 = 0
	if resource.BandwidthBalance < int64(totalBytes) {
		extraBW := int64(totalBytes)
		burnTrx := (float64(extraBW) * 1000) / 1000000
		fmt.Println("bandwidth balance", burnTrx)
		if accountActivated {
			totalTrxNeeded += burnTrx
		}
	}

	energyRequierd, err := EstimateTransactionEnergy(conn, clientAccAddress, contract, merchantAccAddress, balance)
	if err != nil {
		fmt.Println("error estimating transaction energy", err)
		return
	}
	if resource.EnergyBalance < energyRequierd {
		burnTrx := (float64(energyRequierd) * 210) / 1000000
		totalTrxNeeded += burnTrx
		fmt.Println("energy balance", burnTrx)
	}
	// return
	totalTrxNeeded = math.Round(totalTrxNeeded*100000) / 100000
	fmt.Println("totalTrxNeeded", totalTrxNeeded, userAccBalance)

	if userAccBalance < totalTrxNeeded {
		remainigBal := totalTrxNeeded - userAccBalance
		remainigBal = math.Round(remainigBal*100000) / 100000
		result, err := SendTrx(conn, clientAccAddress, merchantAccAddress, merchantAccPrivate, remainigBal)
		if err != nil {
			fmt.Println("error sending trx", err)
			return
		}
		if result == "failed" {
			fmt.Println("failed to send trx")
			return
		}
		fmt.Println("trx send to the client address")

		updateStmt := `UPDATE addresses SET TrxTimeStamp = strftime('%s', 'now') WHERE id = ?;`
		res, err := db.Exec(updateStmt, id) // Replace 1 with the actual id
		if err != nil {
			// log.Fatal(err)
			fmt.Println("error updating the trx time stamp", err)
			return
		}
		fmt.Println("result", res)
	}

	var trxTimeStamp int64
	queryStmt := `SELECT TrxTimeStamp FROM addresses WHERE id = ?;`
	err = db.QueryRow(queryStmt, id).Scan(&trxTimeStamp)
	if err != nil {
		fmt.Println("error getting the trx time stamp", err)
		return
	}

	broadCastTransaction := false
	if trxTimeStamp != 0 {
		fmt.Println("trx", trxTimeStamp)
		currentTime := time.Now().Unix()
		timeDifference := currentTime - trxTimeStamp
		if timeDifference > minSecondsDiff {
			broadCastTransaction = true
		}
	} else {
		if userAccBalance == totalTrxNeeded {
			broadCastTransaction = true
		}
	}

	if broadCastTransaction {
		broadCastResult, err := BroadcastTransaction(conn, tx)
		if err != nil {
			fmt.Println("error broadcasting transaction", err)
			return
		}

		fmt.Println("Transaction broadcasted successfully", broadCastResult, hex.EncodeToString(tx.Txid))
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}

func calculateBytes(tx *api.TransactionExtention) (int, error) {
	var DATA_HEX_PROTOBUF_EXTRA = 3
	var MAX_RESULT_SIZE_IN_TX = 64
	var A_SIGNATURE = 67

	rawData := tx.GetTransaction().GetRawData()
	signatureList := tx.GetTransaction().GetSignature()

	rawDataBytes, err := proto.Marshal(rawData)
	if err != nil {
		fmt.Println("Failed to serialize raw data: ", err)
		return 0, err
	}

	// Calculate base length
	length := len(rawDataBytes) + DATA_HEX_PROTOBUF_EXTRA + MAX_RESULT_SIZE_IN_TX

	// Add signature sizes
	for range signatureList {
		length += A_SIGNATURE
	}
	return length, nil
}
