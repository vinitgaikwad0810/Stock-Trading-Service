// client.go
package main

import (
"fmt"
"log"
"net"
"net/rpc/jsonrpc"
"os"
"strconv"
)


	type Buying_stocks_input struct {
		StockSymbolAndPercentage string
		Budget float64
	}
	type Response_bought_stocks struct {
		trackId int
		stocks []string
		unvestedAmount float64
	}

	type Response_output struct {
		Message string
	}

	func main() {

		argCount := len(os.Args[1:])


		if argCount == 2 {
			fmt.Println("You have opted to buy stocks (PART 1)")

			

			client, err := net.Dial("tcp", "127.0.0.1:1234")
			if err != nil {
				log.Fatal("dialing:", err)
			}

			f, err := strconv.ParseFloat(os.Args[2], 64)
			if err != nil {
				log.Fatal("dialing:", err)
			}
			args := &Buying_stocks_input{os.Args[1],f}
			var reply Response_output
			c := jsonrpc.NewClient(client)
			err = c.Call("Stocktradingsystem.Buyingstocks", args, &reply)
			if err != nil {
				log.Fatal("arith error:", err)
			}
			fmt.Println("\n",reply.Message)
		} else if argCount == 1 {
			client, err := net.Dial("tcp", "127.0.0.1:1234")
			fmt.Println("You have opted for portfolio (PART 2)")
			//fmt.Println("Arguement 1 =" + os.Args[1])

			trackid,err := strconv.Atoi(os.Args[1])
		
			var reply Response_output
			c := jsonrpc.NewClient(client)
			err = c.Call("Stocktradingsystem.CheckingPortfolio", trackid, &reply)
			if err != nil {
				log.Fatal("arith error:", err)
			}
			fmt.Println("\n",reply.Message)

		}else{
			fmt.Println("Incorrect Number of Arguments")
			fmt.Println("Sample Input for Buyingstocks : go run rpc_client.go YHOO:50%,YHOO:30%,GOOG:10% 10000")
		}




	}
