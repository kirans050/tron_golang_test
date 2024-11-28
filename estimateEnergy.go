package main

import (
	"fmt"
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
)

type ResourceBalanceStruct struct {
	BandwidthBalance int64 `json:"bandWidthBalance"`
	EnergyBalance    int64 `json:"energyBalance"`
	FreeNetLimit     int64 `json:"freeNetLimit"`
}

func GetAccountBalance(conn *client.GrpcClient, clientAccAddress string) (float64, error) {
	details, err := conn.GetAccount(clientAccAddress)
	if err != nil {
		return 0, err
	}
	result := float64(details.Balance) / 1000000
	return result, nil
}

func GetAccountResourceHandler(conn *client.GrpcClient, clientAccAddress string) (ResourceBalanceStruct, error) {
	resource, err := conn.GetAccountResource(clientAccAddress)
	if err != nil {
		fmt.Println("error getting resource", err)
		return ResourceBalanceStruct{}, err
	}
	bandwidthBalance := resource.FreeNetLimit - resource.FreeNetUsed
	EnergyBalance := resource.EnergyLimit - resource.EnergyUsed
	fmt.Println("resource", bandwidthBalance, EnergyBalance)
	return ResourceBalanceStruct{
		BandwidthBalance: bandwidthBalance,
		EnergyBalance:    EnergyBalance,
		FreeNetLimit:     resource.FreeNetLimit,
	}, nil
}

func EstimateTransactionEnergy(conn *client.GrpcClient, clientAccAddress, contract, merchantAccAddress string) (*api.EstimateEnergyMessage, error) {

	jsonString := fmt.Sprintf(`[{
		"address":"%s"
	},{
		"uint256":"%s"
	}]`, merchantAccAddress, big.NewInt(10))
	resourceEstimate, err := conn.EstimateEnergy(
		clientAccAddress,
		contract,
		"transfer(address,uint256)",
		jsonString,
		0,
		"",
		0,
	)

	if err != nil {
		fmt.Println("error estimating energy", err)
		return nil, err
	}
	fmt.Println("resource estimate", resourceEstimate)
	return resourceEstimate, nil
}
