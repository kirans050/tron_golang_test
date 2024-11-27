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
}

func GetAccountResourceHandler(conn *client.GrpcClient, fromAddress string) (ResourceBalanceStruct, error) {
	resource, err := conn.GetAccountResource(fromAddress)
	if err != nil {
		fmt.Println("error getting resource", err)
		return ResourceBalanceStruct{}, err
	}
	fmt.Println("resource", resource)
	bandwidthBalance := resource.FreeNetLimit - resource.FreeNetUsed
	EnergyBalance := resource.EnergyLimit - resource.EnergyUsed
	fmt.Println("resource", bandwidthBalance, EnergyBalance)
	return ResourceBalanceStruct{
		BandwidthBalance: bandwidthBalance,
		EnergyBalance:    EnergyBalance,
	}, nil
}

func EstimateTransactionEnergy(conn *client.GrpcClient, fromAddress, contract, toAddress string) (*api.EstimateEnergyMessage, error) {

	jsonString := fmt.Sprintf(`[{
		"address":"%s"
	},{
		"uint256":"%s"
	}]`, toAddress, big.NewInt(10))
	resourceEstimate, err := conn.EstimateEnergy(
		fromAddress,
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
