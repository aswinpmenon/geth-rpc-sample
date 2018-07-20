//create new account in go by JSON RPC
package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
)


func main(){
	client,err :=rpc.Dial("http://localhost:6001")
	if err!=nil{
		fmt.Println("rpc.Dial err",err)
		return
	}
	//declare a string variable to store new generated one
	var account string
	// arg 3 : pwd
	err=client.Call(&account,"personal_newAccount","data")

	if err!=nil{
		fmt.Println("Client.Call err",err)
		return
	}

	fmt.Println(account)

}

