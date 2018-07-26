//Create single raw transactions  with Go-ethereum
package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"io/ioutil"
	"log"
	"math/big"
	"time"
)

func main() {
	// Make a new client based RPC connection to a remote Ethereum node
	client, err := ethclient.Dial("http://localhost:6001") // you can customize
	if err != nil {
		log.Fatalf("Failed to  connetect to the Ethereum client:%v", err)
		return
	}

	// Dir to host the json version of Key
	KeystoreDir := "/Users/data/go/bin/geth-data0/keystore/UTC--2018-07-18T07-41-03.082337980Z--67d7e35ceb75b77fe9f07629f210340c4f53fa57"

	d := time.Now().Add(1000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	keyjson, err := ioutil.ReadFile(KeystoreDir)
	if err != nil {
		log.Fatal(err)
	}

	password := "data"
	// Return a Key instance, private key is one of the fields
	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		log.Fatal("Json key decrypted with bad password")
	}

	// when blockNumber was set nil, its value equals latest by default
	nonce, _ := client.NonceAt(ctx, key.Address, nil)
	fmt.Printf("The nonce of the sending account is %d\n", nonce)

	// Desination address
	receipt := common.HexToAddress("0x20644bb777b9f90c99412fc0cfbaf4b5767c3a26")
	// Value to be sent
	amount := big.NewInt(1000000000)

	// Create a raw tx
	tx := types.NewTransaction(nonce, receipt, amount, big.NewInt(4712388), big.NewInt(0), nil)

	// Make a signature on tx
	tx, err = types.SignTx(tx, types.NewEIP155Signer(big.NewInt(15)), key.PrivateKey)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Derive the sender's address according to tx
	addr, err := types.Sender(types.NewEIP155Signer(big.NewInt(15)), tx)
	if err != nil {
		log.Fatal("Wrong Chain ID")
	}

	// Send a tx
	err = client.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatalf("Tx Error: %s,sender'address:%x,receipt'address:%x,amount sent:%d\n", err, addr, receipt, amount.Div(amount, big.NewInt(params.Ether)))
	} else {
		select {
		case <-time.After(1 * time.Millisecond):
			fmt.Println("overslept")
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		default:
			fmt.Println(tx.Hash().String()) //Print the txhash as a string
		}
	}

}
