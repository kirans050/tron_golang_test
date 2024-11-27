package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/protobuf/proto"
)

func generateSignature(tx *api.TransactionExtention, privateKey string) (*api.TransactionExtention, error) {
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		fmt.Println("error getting rowdata", err)
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	privateKeyBytes, _ := hex.DecodeString(privateKey)
	sk, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	signature, err := crypto.Sign(hash, sk.ToECDSA())
	if err != nil {
		fmt.Println("error signing", err)
		return nil, err
	}
	tx.Transaction.Signature = append(tx.Transaction.Signature, signature)
	return tx, nil
}
