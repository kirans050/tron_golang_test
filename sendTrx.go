package main

import (
	"encoding/hex"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
)

func SendTrx(conn *client.GrpcClient, toAddress, fromAddress, privateKey string) {
	tx, err := conn.Transfer(fromAddress, toAddress, 5)
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
