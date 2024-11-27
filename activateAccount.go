package main

import (
	"encoding/hex"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
)

func ActivateNewAccount(conn *client.GrpcClient, clientAccAddress, merchantAccAddress, merchantAccPrivate string) {
	tx, err := conn.CreateAccount(merchantAccAddress, clientAccAddress)
	if err != nil {
		fmt.Println("error activating account", err)
		return
	}

	tx, err = generateSignature(tx, merchantAccPrivate)
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
