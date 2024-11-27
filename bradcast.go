package main

import (
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
)

func BroadcastTransaction(conn *client.GrpcClient, tx *api.TransactionExtention) (*api.Return, error) {
	result, err := conn.Broadcast(tx.Transaction)
	if err != nil {
		fmt.Println("error broadcasting", err)
		return nil, err
	}
	return result, nil
}
