package main

import (
	"encoding/hex"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
)

type BroadcastStruct struct {
	Result bool
	TxId   string
}

func SendTrx(conn *client.GrpcClient, clientAccAddress, merchantAccAddress, merchantAccPrivate string, totalTrxNeeded float64) (string, error) {
	tx, err := conn.Transfer(merchantAccAddress, clientAccAddress, int64(totalTrxNeeded*1000000))
	if err != nil {
		fmt.Println("error creating transfer", err)
		return "", err
	}

	tx, err = generateSignature(tx, merchantAccPrivate)
	if err != nil {
		fmt.Println("error generating signature", err)
		return "", err
	}

	broadCastResult, err := BroadcastTransaction(conn, tx)
	if err != nil {
		fmt.Println("error broadcasting transaction", err)
		return "", err
	}

	fmt.Println("Transaction broadcasted successfully", broadCastResult, hex.EncodeToString(tx.Txid))

	transactionResult, err := conn.GetTransactionByID(hex.EncodeToString(tx.Txid))
	if err != nil {
		fmt.Println("error getting transaction by id", err)
		return "", err
	}
	val := transactionResult.Ret[0].ContractRet
	if val == core.Transaction_Result_SUCCESS {
		return "success", nil
	} else {
		return "failed", nil
	}
}
