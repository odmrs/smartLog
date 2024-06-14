package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	rpcUrl           = ""
	brlaTokenAddress = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
)

func main() {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("[ERR!] => failed to connect with RPC, error: %v", err)
	}
	defer client.Close()
	log.Println("[INFO] => connection on RPC with success")

	contractAddress := common.HexToAddress(brlaTokenAddress)

	eventSignatureBytes := []byte("Transfer(address,address,uint256)")
	eventSignatureHash := crypto.Keccak256Hash(eventSignatureBytes)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics: [][]common.Hash{
			{eventSignatureHash},
		},
	}

	log.Println("[INFO] => Listening...")

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Printf("========= Transfer Event =========\n"+
				"\tFrom:  %v\n"+
				"\tTo:    %v\n"+
				"-----------------------------------\n", vLog.Topics[1], vLog.Topics[2])
			log.Println("Listening...")
		}
	}
}
