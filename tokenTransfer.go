package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/protobuf/proto"
)

func TokenTransfer(conn *client.GrpcClient, fromAddress, contract, toAddress, privateKey string) {
	balance, err := conn.TRC20ContractBalance(fromAddress, contract)
	if err != nil {
		fmt.Println("error checking token balance", err)
		return
	}
	fmt.Println("balance", balance)
	if balance.Cmp(big.NewInt(0)) <= 0 {
		fmt.Println("insufficent token balance")
		// return
	}

	var bal = big.NewInt(10)
	tx, err := conn.TRC20Send(fromAddress, toAddress, contract, bal, 10000000)
	if err != nil {
		fmt.Println("error generating transaction", err)
		return
	}

	tx, err = generateSignature(tx, privateKey)
	if err != nil {
		fmt.Println("error generating signature", err)
		return
	}
	totalBytes := calculateBytes(tx)
	fmt.Println("totolBytes", totalBytes)

	resource, err := GetAccountResourceHandler(conn, fromAddress)
	if err != nil {
		fmt.Println("unable to get the account resource")
		return
	}

	var totalTrxNeeded float32 = 0
	if resource.BandwidthBalance < int64(totalBytes) {
		extraBW := int64(totalBytes)
		// extraBW := int64(totalBytes) - resource.BandwidthBalance
		burnTrx := (float32(extraBW) * 1000) / 1000000
		totalTrxNeeded += burnTrx
	}

	result, err := EstimateTransactionEnergy(conn, fromAddress, contract, toAddress)
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
		burnTrx := (float32(13091) * 210) / 1000000
		totalTrxNeeded += burnTrx
	}
	// return
	fmt.Println("totalTrxNeeded", totalTrxNeeded)

	// broadCastResult, err := BroadcastTransaction(conn, tx)
	// if err != nil {
	// 	fmt.Println("error broadcasting transaction", err)
	// 	return
	// }

	// fmt.Println("Transaction broadcasted successfully", broadCastResult, hex.EncodeToString(tx.Txid))
}

func calculateBytes(tx *api.TransactionExtention) int {
	var DATA_HEX_PROTOBUF_EXTRA = 3
	var MAX_RESULT_SIZE_IN_TX = 64
	var A_SIGNATURE = 67

	rawData := tx.GetTransaction().GetRawData()
	signatureList := tx.GetTransaction().GetSignature()

	rawDataBytes, err := proto.Marshal(rawData)
	if err != nil {
		log.Fatalf("Failed to serialize raw data: %v", err)
	}

	// Calculate base length
	length := len(rawDataBytes) + DATA_HEX_PROTOBUF_EXTRA + MAX_RESULT_SIZE_IN_TX

	// Add signature sizes
	for range signatureList {
		length += A_SIGNATURE
	}

	fmt.Println("lenght", length)
	return length
}
