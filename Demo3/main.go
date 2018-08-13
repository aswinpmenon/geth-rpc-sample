//Create batches of transactions  with Go-ethereum

package main

import (
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"strings"
	"time"
	"bufio"
	"io"
	"sync"
)

const ACC_NUM = 4

var mux sync.Mutex

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()
	client, err := ethclient.Dial("http://localhost:6001")
	if err != nil {
		fmt.Println("rpc.Dial err", err)
		return
	}

	//fileName :="accounts.txt"

	//dstFile, err:= os.Create(fileName)  //usr/data/go/accounts.txt
	//if err!=nil{
	//	fmt.Println(err.Error())
	//	return
	//}
	//defer dstFile.Close()

	//declare a string variable to store new generated one
	//var account string
	// arg 3 : pwd

	// Create ACC_NUM accounts
	//for i:=0;i<ACC_NUM ;i++{
	//	err = client.Call(&account, "personal_newAccount", "data")
	//	if err != nil {
	//		fmt.Println("Client.Call err", err)
	//		return
	//	}
	//	// write account into the created file
	//	dstFile.WriteString(account+"\n")
	//}

	//Scan the accounts within given file

	var accounts []string

	accountFile := "/Users/data/go/accounts.txt"

	src_file, inputError := os.Open(accountFile)
	if inputError != nil {
		fmt.Println(inputError)
		return
	}
	defer src_file.Close()


	// read file by row
	buf := bufio.NewReader(src_file)
	for {
		line, err := buf.ReadString('\n')

		//fmt.Println(line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
		accounts= append(accounts,strings.TrimSpace(line))
	}

	//for testing
	for _, value := range accounts {
		fmt.Println(value)
	}

	// Decrypt them
	var fiName []string

	dir_dst, err := ioutil.ReadDir("/Users/data/go/bin/geth-data0/keystore")
	if err != nil {
		log.Fatalf("read dir error:", err)
		return
	}

	for _, fi := range dir_dst {
		//The second clause is optional,Whether needs it depends on your real dir
		if fi.IsDir() || strings.HasPrefix(fi.Name(), "package-lock.json") {
			continue
		} else {
			fiName = append(fiName, fi.Name()) //only save file
		}
	}

	//So far, so good
	//for testing
	for _, file := range fiName {
		fmt.Println(file)
	}

	var num = len(fiName)
	var pwd = "data"
	exitChan := make(chan int)

	amount:=big.NewInt(1)

	// the first
	for i := 0; i < num; i++ {
		go func(fileIndex string, pwd string, amount *big.Int, client *ethclient.Client,ctx context.Context, fiName []string ) {
			mux.Lock()
			defer mux.Unlock()
			keyjson, err := ioutil.ReadFile(fileIndex)
			if err != nil {
				log.Fatal(err)
			}
			key, err := keystore.DecryptKey(keyjson, pwd)

			nonce, _ := client.NonceAt(ctx, key.Address, nil)
			fmt.Printf("The nonce of the sending account is %d\n", nonce) //the nonce is out of order due to concurrent
			for j:=0;j<len(fiName);j++{
				dst:=common.HexToAddress(fiName[j])
				txHash,txContent:=makeRawTransaction(&dst,amount,client,nonce,ctx,key)
				fmt.Printf("tx in hash:%s its content is %s\n",txHash,txContent)
			}
			exitChan <- 1
		}("/Users/data/go/bin/geth-data0/keystore/"+fiName[i],pwd,amount,client,ctx,fiName) //add the common prefix of each keyfile

	}

	for i := 0; i < num; i++ {
		<-exitChan
	}

	//fmt.Println(account)

}

//const ROUTING_NUM=10
//
//var mutex sync.Mutex
//
//func main() {
//	//Make a new client based RPC connection to a remote Ethereum node
//	client, err := ethclient.Dial("http://localhost:6001")  // you can customize
//	if err != nil {
//		log.Fatalf("Failed to  connetect to the Ethereum client:%v", err)
//		return
//	}
//
//	var accounts [2000]string
//
//
//	// Dir to host the json version of Key
//	KeystoreDir := "/Users/data/go/bin/geth-data0/keystore/UTC--2018-07-18T07-41-03.082337980Z--67d7e35ceb75b77fe9f07629f210340c4f53fa57"
//
//	//d:=time.Now().Add(1000*time.Millisecond)
//	d:=time.Now().Add(30*time.Second)
//	ctx,cancel:=context.WithDeadline(context.Background(),d)
//	defer cancel()  //30s later, execute cancel()
//
//	keyjson, err := ioutil.ReadFile(KeystoreDir)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	password := "data"
//	// Return a Key instance, private key is one of the fields
//	key, err := keystore.DecryptKey(keyjson, password)
//	if err != nil {
//		log.Fatal("Json key decrypted with bad password")
//	}
//
//	// when blockNumber was set nil, its value equals latest by default
//	nonce,_:=client.NonceAt(ctx,key.Address,nil)
//	fmt.Printf("The nonce of the sending account is %d\n",nonce)
//
//	// Desination address
//	receipt:=common.HexToAddress("0x20644bb777b9f90c99412fc0cfbaf4b5767c3a26")
//	// Value to be sent
//	amount:=big.NewInt(1000000000)
//
//	txNum,_:=strconv.Atoi(os.Args[1])
//	txNum=txNum*10000
//
//	if txNum> 2<<32{
//		txNum=2<<32
//	}
//
//	exitChain:=make(chan int)
//	txNumPerRoutine:=txNum/ROUTING_NUM
//
//	counter:=0
//
//	for i:=1;i<=ROUTING_NUM;i++{
//		go func(nonce uint64,routineIndex int) {
//			mutex.Lock()
//			for j:=0;j<txNumPerRoutine;j++{
//
//				txHash,txContent:=makeRawTransaction(&receipt,amount,client,nonce,ctx,key)
//				nonce++
//				counter++
//
//				fmt.Println(txHash,",",txContent)
//			}
//			exitChain<-1
//			mutex.Unlock()
//		}(nonce,i)
//	}
//
//	for i:=0;i<ROUTING_NUM;i++{
//		<-exitChain
//	}
//
//	fmt.Println("计数器：",counter)
//
//
//
//}
//
func makeRawTransaction(to *common.Address, amount *big.Int, client *ethclient.Client, nonce uint64, ctx context.Context, key *keystore.Key) (string, string) {

	// Create a raw tx
	tx := types.NewTransaction(nonce, *to, amount, big.NewInt(21000), big.NewInt(0), nil)

	// Make a signature on tx
	tx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(15)), key.PrivateKey)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}

	// Derive the sender's address according to tx
	addr, err := types.Sender(types.NewEIP155Signer(big.NewInt(15)), tx)
	if err != nil {
		log.Fatal("Wrong Chain ID")
		return "", ""
	}

	// Send a tx
	err = client.SendTransaction(ctx, tx)

	var txbf []byte
	if txbf, err = tx.MarshalJSON(); err != nil {
		fmt.Println("Serialization of transaction tx err")
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("Tx Error: %s,sender'address:%x,receipt'address:%x,amount sent:%d\n", err, addr, to, amount.Div(amount, big.NewInt(params.Ether)))
		return "", ""
	} else {
		select {
		case <-time.After(1 * time.Millisecond):
			fmt.Println("overslept")
		case <-ctx.Done():
			fmt.Printf("err:%s\n", ctx.Err())

		default:
			//Print the txhash as a string
			return tx.Hash().String(), common.ToHex([]byte(txbf))
		}
	}
	return "", ""

}

