package main

import (
	"encoding/hex"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
)

func SendTrx(conn *client.GrpcClient, clientAccAddress, fromAddress, privateKey string, totalTrxNeeded float32) {
	fmt.Println("int64(totalTrxNeeded)", int64(totalTrxNeeded), totalTrxNeeded)
	tx, err := conn.Transfer(fromAddress, clientAccAddress, int64(totalTrxNeeded*1000000))
	if err != nil {
		fmt.Println("error creating transfer", err)
		return
	}

	tx, err = generateSignature(tx, privateKey)
	if err != nil {
		fmt.Println("error generating signature", err)
		return
	}

	broadCastResult, err := BroadcastTransaction(conn, tx)
	if err != nil {
		fmt.Println("error broadcasting transaction", err)
		return
	}

	fmt.Println("Transaction broadcasted successfully", broadCastResult, hex.EncodeToString(tx.Txid))
}
