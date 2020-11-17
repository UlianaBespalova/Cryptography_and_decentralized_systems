package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var PRIV_KEY = "e32229e91ae57ad4baf08e2d508e45e0bd4a27440f415427a7bde58d67b3aaa8"

var GAS_PRICE = 3e+10
var GAS_LIMIT = 21000

var (
	toFlag = flag.String("to", "", "recipient address")
	valueFlag = flag.String("value", "", "transfer amount (ETH)")
)

type GethTransaction struct {
	To   string     `json:"to"`
	From string     `json:"from"`
	Nonce string    `json:"nonce"`
	Gas string      `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value string    `json:"value"`
	Data string     `json:"input"`
	V string        `json:"v"`
	R string     	`json:"r"`
	S string        `json:"s"`
	Hash string      `json:"hash"`
}


func EthToWai (valueEth *big.Float) (*big.Int) {
	var EthToWaiK = big.NewFloat(1e+18)
	value := new(big.Float)
	value.Mul(valueEth, EthToWaiK)
	valueWai := new(big.Int)
	value.Int(valueWai)
	return valueWai
}

func generateTransaction (toStr string, valueStr string) {

	privateKey, err := crypto.HexToECDSA(PRIV_KEY)
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("Error: failed to cast public key to ECDSA")
		return
	}
	addressFrom := crypto.PubkeyToAddress(*publicKeyECDSA)
	addressTo := common.HexToAddress(toStr)

	valueEth := new(big.Float)
	valueEth, ok = valueEth.SetString(valueStr)
	if !ok {
		fmt.Println("Error: invalid number of ethers")
		return
	}
	valueWai := EthToWai(valueEth)

	gasLimit := uint64(GAS_LIMIT)
	gasPrice := big.NewInt(int64(GAS_PRICE))


	var nonce uint64 = 0
	tx := types.NewTransaction(nonce, addressTo, valueWai, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(nil), privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	var parsedTx = new(GethTransaction)
	jsonTx, _ := signedTx.MarshalJSON()
	_ = json.Unmarshal(jsonTx, &parsedTx) //add field "from"
	parsedTx.From = addressFrom.String()
	parsedTxJson, _ := json.Marshal(parsedTx)

	fmt.Println(string(parsedTxJson))
}


func main() {

	flag.Parse()
	if *toFlag!="" && *valueFlag!="" {
		generateTransaction(*toFlag, *valueFlag)
		return
	}
	fmt.Println("Flags error!")
	return
}
