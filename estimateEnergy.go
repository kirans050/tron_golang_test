package main

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
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

func EstimateTransactionEnergy(conn *client.GrpcClient, clientAccAddress, contract, merchantAccAddress string, balance *big.Int) (int64, error) {

	jsonString := fmt.Sprintf(`[{
		"address":"%s"
	},{
		"uint256":"%s"
	}]`, merchantAccAddress, balance)

	// resourceEstimate, err := conn.EstimateEnergy(
	// 	clientAccAddress,
	// 	contract,
	// 	"transfer(address,uint256)",
	// 	jsonString,
	// 	0,
	// 	"",
	// 	0,
	// )

	// if err != nil {
	// 	fmt.Println("error estimating energy", err)
	// 	return nil, err
	// }

	val, err := conn.TriggerConstantContract(clientAccAddress, contract, "transfer(address,uint256)", jsonString)
	if err != nil {
		fmt.Println("error triggering contract", err)
		return 0, err
	} else {
		if !val.Result.Result {
			return 0, errors.New("estimate failed")
		}
		return val.EnergyUsed, nil
	}

	// fmt.Println("resource estimate", resourceEstimate)
	// return resourceEstimate, nil
}
