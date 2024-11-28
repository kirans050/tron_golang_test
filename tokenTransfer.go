package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/protobuf/proto"
)

func TokenTransfer(conn *client.GrpcClient, clientAccAddress, contract, merchantAccAddress, clientAccPrivate, merchantAccPrivate string) {
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
	accBalance, err := GetAccountBalance(conn, clientAccAddress)
	if err != nil {
		if err.Error() == "account not found" {
			userAccBalance = 0
		} else {
			fmt.Println("error getting balance", err)
			return
		}
	} else {
		userAccBalance = float64(accBalance)
	}
	fmt.Println("acc", accBalance)

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
		totalTrxNeeded += burnTrx
		fmt.Println("bandwidth balance", burnTrx)
	}

	result, err := EstimateTransactionEnergy(conn, clientAccAddress, contract, merchantAccAddress)
	if err != nil {
		fmt.Println("error estimating transaction energy", err)
		return
	}
	fmt.Println("result", result)
	if !result.Result.Result {
		fmt.Println("unable to fetch transaction energy", err)
		return
	}
	energyRequierd := result.EnergyRequired
	if resource.EnergyBalance < energyRequierd {
		burnTrx := (float64(13091) * 210) / 1000000
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
	}

	broadCastResult, err := BroadcastTransaction(conn, tx)
	if err != nil {
		fmt.Println("error broadcasting transaction", err)
		return
	}

	fmt.Println("Transaction broadcasted successfully", broadCastResult, hex.EncodeToString(tx.Txid))
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
