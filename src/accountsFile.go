package main

import (
	"fmt"
	"os"
	"github.com/ethererum/go-ethereum/ethclient"
)

const ACC_NUM = 4  //the number of accounts

func main() {
	// Make a connection to geth
	client, err := ethclient.Dial("http://localhost:6001")
	if err != nil {
		fmt.Println("rpc.Dial err", err)
		return
	}

	//save new accounts
	fileName := "accounts.txt"

	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()

	//declare a string variable to store new generated one
	var account string

	//Create ACC_NUM accounts
	for i := 0; i < ACC_NUM; i++ {
		// 3rd arg is pwd
		err = client.Call(&account, "personal_newAccount", "data")
		if err != nil {
			fmt.Println("Client.Call err", err)
			return
		}
		// write account into the created file
		dstFile.WriteString(account + "\n")
	}

}
