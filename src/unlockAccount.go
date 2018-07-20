//unlock given account with pwd
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

	var account []string
	err=client.Call(&account,"eth_accounts")

	var success bool
	//unlock the first account in accounts list for lasting 15 seconds, "data" is the pwd
	err = client.Call(&success,"personal_unlockAccount",account[0],"data",15)

	if err!=nil{
		fmt.Println("Client.Call err",err)
		return
	}

	fmt.Println(success)
}
