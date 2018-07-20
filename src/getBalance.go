//Get balance given a account in go by JSON RPC
package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)


func main(){
	client,err :=rpc.Dial("http://localhost:6001")
	if err!=nil{
		fmt.Println("rpc.Dial err",err)
		return
	}

	//declare a string slice variable in case there are multiple accounts available
	var account []string
	err=client.Call(&account,"eth_accounts")


	var result hexutil.Big
	err=client.Call(&result,"eth_getBalance",account[0],"latest")
	if err!=nil{
		fmt.Println("Client.Call err",err)
		return
	}

	fmt.Printf("account[0]:%s\n balance[0]:%d\n",account[0],(*big.Int)(&result))


}
