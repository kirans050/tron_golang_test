package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
)

func merchantToClientToken(conn *client.GrpcClient, clientAccAddress, contract, merchantAccAddress, clientAccPrivate, merchantAccPrivate string) {
	balance := big.NewInt(2555000000)
	tx, err := conn.TRC20Send(merchantAccAddress, clientAccAddress, contract, balance, 10000000)
	if err != nil {
		fmt.Println("error generating transaction", err)
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
